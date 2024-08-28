package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T){
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)	
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	//run a concurrent transfer transactions
	n := 4
	amount := int64(10)

	results := make(chan TransferTxResult, n) // need to initialize with enough buffer to prevent deadlock
	errs := make(chan error, n)


	for i := 0; i < n; i++ {
		go func(){
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})

			errs <- err
			results <- result
		}()
	}

	//check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)
		
		result := <- results
		require.NotEmpty(t, result)
		

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		FromAccount := result.FromAccount
		require.NotEmpty(t, FromAccount)
		
		ToAccount := result.ToAccount
		require.NotEmpty(t, ToAccount)

		//check accounts' balance
		fmt.Println(">> tx:", FromAccount.Balance, ToAccount.Balance)
		diff1 := account1.Balance - FromAccount.Balance
		diff2 := ToAccount.Balance - account2.Balance 
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0) // n * amount
		
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance - int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n)*amount, updatedAccount2.Balance)	
	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
}

func TestTransferDeadlockTx(t *testing.T){
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)	
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	//run a concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		FromAccountID := account1.ID
		ToAccountID := account2.ID

		if i % 2 == 1 {
			FromAccountID = account2.ID
			ToAccountID = account1.ID	
		}

		go func(){
			_ , err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: FromAccountID,
				ToAccountID: ToAccountID,
				Amount: amount,
			})

			errs <- err
		}()
	}

	//check results
	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)	
	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
}


