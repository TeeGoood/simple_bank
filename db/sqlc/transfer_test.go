package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/teegoood/simplebank/util"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: util.RandomMoney(),
	}

	tranfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tranfer)

	require.Equal(t, tranfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, tranfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, tranfer.Amount, arg.Amount)	
	
	require.NotZero(t, tranfer.ID)
	require.NotZero(t, tranfer.CreatedAt)

	return tranfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T){
	transfer1 := createRandomTransfer(t)
	
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)	
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		arg := CreateTransferParams{
			FromAccountID: account1.ID,
			ToAccountID: account2.ID,
			Amount: util.RandomMoney(),
		}

		testQueries.CreateTransfer(context.Background(), arg)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Offset: 5,
		Limit: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, tranfer := range transfers {
		require.NotEmpty(t, tranfer)
		require.Equal(t, tranfer.FromAccountID, account1.ID)
		require.Equal(t, tranfer.ToAccountID, account2.ID)
	} 
}