package wallet

import (
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	type args struct {
		seed string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"1st", args{""}, nil},
		{"2nd", args{"12345678912345678912345678912345"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (<-GenerateKey(tt.args.seed)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreate_PathSet(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	home := currentUser.HomeDir
	path := filepath.Join(home, "/.indy_client/")

	sc := &StorageConfig{Path: path}

	type args struct {
		config      Config
		credentials Credentials
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"1st", args{Config{ID: "test1_wallet_in_root", StorageConfig: sc},
			Credentials{Key: "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp", KeyDerivationMethod: "RAW"}}, nil},
		{"2nd", args{Config{ID: "test2_wallet_in_root", StorageConfig: sc},
			Credentials{Key: "C7mR5TZVB7WRCYsTMQGXuLHcXisFYZL1GoXARyiVyEER", KeyDerivationMethod: "RAW"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (<-Create(tt.args.config, tt.args.credentials)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
			walletPath := filepath.Join(path, tt.args.config.ID)
			if _, err := os.Stat(walletPath); os.IsNotExist(err) {
				t.Errorf("wallet path (%s) does not exist", walletPath)
			}
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (<-Delete(tt.args.config, tt.args.credentials)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		config      Config
		credentials Credentials
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"1st", args{Config{ID: "test1_wallet"}, Credentials{Key: "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp", KeyDerivationMethod: "RAW"}}, nil},
		{"2nd", args{Config{ID: "test2_wallet"}, Credentials{Key: "C7mR5TZVB7WRCYsTMQGXuLHcXisFYZL1GoXARyiVyEER", KeyDerivationMethod: "RAW"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (<-Create(tt.args.config, tt.args.credentials)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenAndClose(t *testing.T) {
	type args struct {
		config      Config
		credentials Credentials
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"1st", args{Config{ID: "test1_wallet"}, Credentials{Key: "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp", KeyDerivationMethod: "RAW"}}, nil},
		{"2nd", args{Config{ID: "test2_wallet"}, Credentials{Key: "C7mR5TZVB7WRCYsTMQGXuLHcXisFYZL1GoXARyiVyEER", KeyDerivationMethod: "RAW"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := <-Open(tt.args.config, tt.args.credentials)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() = %v, want %v", got, tt.want)
			}
			if got := (<-Close(r.Handle())).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Close() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestOpen2Times(t *testing.T) {
//	type args struct {
//		config      Config
//		credentials Credentials
//	}
//	tests := []struct {
//		name string
//		args args
//		want error
//	}{
//		{"1st", args{Config{ID: "unit_test_wallet1"}, Credentials{Key: "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp", KeyDerivationMethod: "RAW"}}, nil},
//	//	{"2nd", args{Config{ID: "test2_wallet"}, Credentials{Key: "C7mR5TZVB7WRCYsTMQGXuLHcXisFYZL1GoXARyiVyEER", KeyDerivationMethod: "RAW"}}, nil},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := <-Open(tt.args.config, tt.args.credentials)
//			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Open() = %v, want %v", got, tt.want)
//			}
//			firstWallet := r.Handle()
//			fmt.Println(firstWallet)
//
//			r = <-Open(tt.args.config, tt.args.credentials)
//			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Open() = %v, want %v", got, tt.want)
//			}
//			fmt.Println(r.Handle())
//
//			if got := (<-Close(r.Handle())).Err(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Close() = %v, want %v", got, tt.want)
//			}
//			if got := (<-Close(firstWallet)).Err(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Close() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
// -------------------------------------------------
// ->RESULT: CANNOT OPEN SAME WALLET MORE THAN ONCE!
// -------------------------------------------------

func TestExport(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	home := currentUser.HomeDir
	tmpfn := filepath.Join(home, "/.indy_client/wallet/test_wallet1.export")

	type args struct {
		config      Config
		credentials Credentials
		exportCfg   Credentials
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"1st",
			args{
				Config{ID: "test1_wallet"},
				Credentials{
					Key:                 "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp",
					KeyDerivationMethod: "RAW"},
				Credentials{
					Path:                tmpfn,
					Key:                 "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp",
					KeyDerivationMethod: "RAW"},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := <-Open(tt.args.config, tt.args.credentials)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() = %v, want %v", got, tt.want)
			}
			w := r.Handle()
			if got := (<-Export(w, tt.args.exportCfg)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Export() = %v, want %v", got, tt.want)
			}
			if got := (<-Close(w)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Close() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImport(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	home := currentUser.HomeDir
	tmpfn := filepath.Join(home, "/.indy_client/wallet/test_wallet1.export")

	type args struct {
		config      Config
		credentials Credentials
		importCfg   Credentials
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"1st",
			args{
				Config{ID: "test3_wallet"},
				Credentials{Key: "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp", KeyDerivationMethod: "RAW"},
				Credentials{Path: tmpfn, Key: "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp"},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (<-Import(tt.args.config, tt.args.credentials, tt.args.importCfg)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Import() = %v, want %v", got, tt.want)
			}
			r := <-Open(tt.args.config, tt.args.credentials)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() = %v, want %v", got, tt.want)
			}
			if got := (<-Close(r.Handle())).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Close() = %v, want %v", got, tt.want)
			}
		})
	}
	_ = os.Remove(tmpfn)
}

func TestDelete(t *testing.T) {
	type args struct {
		config      Config
		credentials Credentials
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"1st", args{Config{ID: "test1_wallet"}, Credentials{Key: "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp", KeyDerivationMethod: "RAW"}}, nil},
		{"2nd", args{Config{ID: "test2_wallet"}, Credentials{Key: "C7mR5TZVB7WRCYsTMQGXuLHcXisFYZL1GoXARyiVyEER", KeyDerivationMethod: "RAW"}}, nil},
		{"exported", args{Config{ID: "test3_wallet"}, Credentials{Key: "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp", KeyDerivationMethod: "RAW"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (<-Delete(tt.args.config, tt.args.credentials)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
