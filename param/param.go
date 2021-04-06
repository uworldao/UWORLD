package param

import "github.com/uworldao/UWORLD/common/hasharry"

const (
	// Mainnet logo
	MainNet = "MN"
	// Testnet logo
	TestNet = "TN"
)

var (
	// Program name
	AppName = "UWorld"
	// Current network
	Net = MainNet
	// Network logo
	UniqueNetWork = "_UWorld"
	// Token name
	Token = hasharry.StringToAddress("UWD")

	FeeAddress   = hasharry.StringToAddress("UWDNQhgkNHCLdVhCFvpo6bGXXdcKtTTfeQZE")
	EaterAddress = hasharry.StringToAddress("UWDCoinEaterAddressDontSend000000000")
)

const (
	// Block interval period
	BlockInterval = uint64(30)
	// Re-election interval
	TermInterval = 60 * 60 * 24 * 365 * 100
	// Maximum number of super nodes
	MaxWinnerSize = 11
	// The minimum number of nodes required to confirm the transaction
	SafeSize = MaxWinnerSize*2/3 + 1
	// The minimum threshold at which a block is valid
	ConsensusSize                 = MaxWinnerSize*2/3 + 1
	SkipCurrentWinnerWaitTimeBase = BlockInterval * MaxWinnerSize * 1
)

const (
	// AtomsPerCoin is the number of atomic units in one coin.
	AtomsPerCoin = 1e8

	// Circulation is the total number of COINS issued.
	Circulation = 6300 * 1e4 * AtomsPerCoin

	// GenesisCoins is genesis Coins
	GenesisCoins = 210 * 1e4 * AtomsPerCoin

	// CoinBaseCoins is reward
	CoinBaseCoins = 3 * AtomsPerCoin

	//MaxAddressTxs is address the maximum number of transactions in the trading pool
	MaxAddressTxs = 1000

	// MinFeesCoefficient is minimum fee required to process the transaction
	MinFeesCoefficient uint64 = 1e4

	// MaxFeesCoefficient is maximum fee required to process the transaction
	MaxFeesCoefficient uint64 = 1 * AtomsPerCoin

	// MinAllowedAmount is the minimum allowable amount for a transaction
	MinAllowedAmount uint64 = 0.005 * AtomsPerCoin

	// MaxAllContractCoin is the maximum allowable sum of contract COINS
	MaxAllContractCoin uint64 = 1e11 * AtomsPerCoin

	// MaxContractCoin is the maximum allowable contract COINS
	MaxContractCoin uint64 = 1e10 * AtomsPerCoin

	Fees uint64 = 0.002 * AtomsPerCoin

	TokenConsumption uint64 = 10.24 * AtomsPerCoin

	CoinHeight = 1
)

var (
	MainPubKeyHashAddrID  = [3]byte{0x03, 0x82, 0x32} //UWD 3, 82, 32
	TestPubKeyHashAddrID  = [3]byte{0x03, 0x82, 0x32} //uwd
	MainPubKeyHashTokenID = [3]byte{0x03, 0x82, 0x55} //UWT 3, 82, 55
	TestPubKeyHashTokenID = [3]byte{0x03, 0x82, 0x55} //uwt
)

type MappingInfo struct {
	Address string
	Note    string
	Amount  uint64
}

var MappingCoin = []MappingInfo{
	{
		Address: "UWDM1qcsk7UUNANMPKSpALJW7AqpDCy7tdoN",
		Note:    "",
		Amount:  210 * 1e4 * 1e8,
	},
}

var DayCoin = map[uint64]float64{
	1:   7000,
	2:   7700,
	3:   8470,
	4:   9317,
	5:   10248.7,
	6:   11273.57,
	7:   9920.7416,
	8:   10714.40093,
	9:   11571.553,
	10:  12497.27724,
	11:  13497.05942,
	12:  14576.82418,
	13:  15742.97011,
	14:  17002.40772,
	15:  18362.60034,
	16:  19831.60836,
	17:  21418.13703,
	18:  23131.58799,
	19:  15613.8219,
	20:  16394.51299,
	21:  17214.23864,
	22:  18074.95057,
	23:  18978.6981,
	24:  19927.63301,
	25:  20924.01466,
	26:  21970.21539,
	27:  23068.72616,
	28:  24222.16247,
	29:  25433.27059,
	30:  26704.93412,
	31:  28040.18083,
	32:  29442.18987,
	33:  30914.29936,
	34:  32460.01433,
	35:  34083.01504,
	36:  35787.1658,
	37:  37576.52409,
	38:  39455.35029,
	39:  41428.11781,
	40:  43499.5237,
	41:  45674.49988,
	42:  47958.22487,
	43:  30213.68167,
	44:  31120.09212,
	45:  32053.69488,
	46:  33015.30573,
	47:  34005.7649,
	48:  35025.93785,
	49:  36076.71599,
	50:  37159.01747,
	51:  38273.78799,
	52:  39422.00163,
	53:  40604.66168,
	54:  41822.80153,
	55:  43077.48557,
	56:  44369.81014,
	57:  45700.90445,
	58:  47071.93158,
	59:  48484.08953,
	60:  49938.61221,
	61:  51436.77058,
	62:  52979.8737,
	63:  54569.26991,
	64:  56206.348,
	65:  57892.53844,
	66:  59629.3146,
	67:  61418.19404,
	68:  63260.73986,
	69:  65158.56205,
	70:  67113.31891,
	71:  69126.71848,
	72:  71200.52004,
	73:  73336.53564,
	74:  75536.63171,
	75:  77802.73066,
	76:  80136.81258,
	77:  82540.91695,
	78:  85017.14446,
	79:  87567.6588,
	80:  90194.68856,
	81:  92900.52922,
	82:  95687.54509,
	83:  98558.17145,
	84:  101514.9166,
	85:  69706.90939,
	86:  71101.04758,
	87:  72523.06853,
	88:  73973.5299,
	89:  75453.0005,
	90:  76962.06051,
	91:  78501.30172,
	92:  80071.32775,
	93:  81672.75431,
	94:  83306.2094,
	95:  84972.33358,
	96:  86671.78026,
	97:  88405.21586,
	98:  90173.32018,
	99:  91976.78658,
	100: 93816.32231,
	101: 95692.64876,
	102: 97606.50173,
	103: 99558.63177,
	104: 101549.8044,
	105: 103580.8005,
	106: 105652.4165,
	107: 107765.4648,
	108: 109920.7741,
	109: 112119.1896,
	110: 114361.5734,
	111: 116648.8049,
	112: 118981.781,
	113: 121361.4166,
	114: 123788.6449,
	115: 126264.4178,
	116: 128789.7062,
	117: 131365.5003,
	118: 133992.8103,
	119: 136672.6665,
	120: 139406.1198,
}
