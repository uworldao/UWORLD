package ut

import (
	"fmt"
	"github.com/uworldao/UWORLD/crypto/ecc/secp256k1"
	"github.com/uworldao/UWORLD/param"
	"testing"
)

func TestGenerateUWDAddress(t *testing.T) {
	key, _ := secp256k1.GeneratePrivateKey()

	address := GenerateUWDAddress(param.Net, key.PubKey())
	fmt.Println(address.String())
	if !CheckUWDAddress(param.Net, address.String()) {
		t.Fatal("failed")
	}
}

func TestCheckCoinName(t *testing.T) {
	type args struct {
		coinName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{"Bit"}, false},
		{"2", args{"BIT"}, false},
		{"3", args{"bit"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckAbbr(tt.args.coinName); (err != nil) != tt.wantErr {
				t.Errorf("CheckCoinName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateContractAddress(t *testing.T) {
	type args struct {
		version string
		address string
		name    string
	}
	key, _ := secp256k1.GeneratePrivateKey()
	address := GenerateUWDAddress("TN", key.PubKey())
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"1", args{"TN", address.String(), "BIT"}, true, false},
		{"2", args{"TN", address.String(), "BIT"}, true, false},
		{"3", args{"TN", address.String(), "BBB"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateContractAddress(tt.args.version, tt.args.address, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateContractAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			check := CheckContractAddress(tt.args.version, tt.args.address, tt.args.name, got.String())
			if check != tt.want {
				t.Errorf("GenerateContractAddress() got = %v, want %v", check, tt.want)
			}
		})
	}
}
