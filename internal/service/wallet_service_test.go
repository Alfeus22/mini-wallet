package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alfeus22/mini-wallet/internal/service"
)

type MockWalletStorage struct {
	mock.Mock
}

func (m *MockWalletStorage) TransferTx(pengirimID int, penerimaID int, jumlah int) error {
	args := m.Called(pengirimID, penerimaID, jumlah)

	return args.Error(0)
}

func TestTransfer(t *testing.T) {
	// Daftar skenario pengujian
	tests := []struct {
		name          string
		pengirimID    int
		penerimaID    int
		amount        int
		mockSetup     func(mockStorage *MockWalletStorage)
		expectedError string
	}{
		{
			name:       "Skenario 1: Transfer Sukses",
			pengirimID: 1,
			penerimaID: 2,
			amount:     50000,
			mockSetup: func(mockStorage *MockWalletStorage) {
				// Kita suruh gudang palsu menjawab "nil" (tidak ada error)
				mockStorage.On("TransferTx", 1, 2, 50000).Return(nil)
			},
			expectedError: "",
		},
		{
			name:       "Skenario 2: Ditolak - Nominal Minus",
			pengirimID: 1,
			penerimaID: 2,
			amount:     -10000,
			mockSetup: func(mockStorage *MockWalletStorage) {
				// Gudang JANGAN dipanggil karena Service harus menolaknya duluan!
			},
			expectedError: "Jumlah transfer harus lebih dari 0",
		},
		{
			name:       "Skenario 3: Ditolak - Transfer ke Diri Sendiri",
			pengirimID: 1,
			penerimaID: 1,
			amount:     50000,
			mockSetup: func(mockStorage *MockWalletStorage) {
				// Gudang JANGAN dipanggil
			},
			expectedError: "transfer yang bener",
		},
		{
			name:       "Skenario 4: Ditolak Gudang - Saldo Kurang",
			pengirimID: 1,
			penerimaID: 2,
			amount:     9999999,
			mockSetup: func(mockStorage *MockWalletStorage) {
				// Kita suruh gudang palsu pura-pura mengecek dan mengembalikan error
				mockStorage.On("TransferTx", 1, 2, 9999999).Return(errors.New("saldo tidak cukup"))
			},
			expectedError: "saldo tidak cukup",
		},
	}

	// Jalankan semua skenario di atas secara looping
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Siapkan gudang palsu yang baru
			mockStorage := new(MockWalletStorage)
			tc.mockSetup(mockStorage)

			// 2. Masukkan gudang palsu ke dapur (Service)
			walletService := service.NewWalletService(mockStorage)

			// 3. Eksekusi transfer
			err := walletService.TransferTx(tc.pengirimID, tc.penerimaID, tc.amount)

			// 4. Validasi hasil (Assert)
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}

			// 5. Pastikan perilaku gudang sesuai dengan mockSetup
			mockStorage.AssertExpectations(t)
		})
	}
}
