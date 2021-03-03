package addons

import (
	"context"
	"os"
	"testing"

	"github.com/codenotary/immudb/pkg/api/schema"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/stretchr/testify/assert"
)

// Setting this as true means that testing is done against real ImmuDB
// instead of mocked immuDB. So this can be used as an E2E testing tool
const testAgainstRealImmuDB = false

const immuTxnIDForClaim = "2Kc7X1ErDwNQC3mDSzcj2r:3:CL:2Kc7X1ErDwNQC3mDSzcj2r:2:NEW_SCHEMA_58906353:1.0:TAG_1"
const immuClaimDataToWrite = `
	{
		"ver":"1.0",
		"id":"2Kc7X1ErDwNQC3mDSzcj2r:3:CL:2Kc7X1ErDwNQC3mDSzcj2r:2:NEW_SCHEMA_58906353:1.0:TAG_1",
		"schemaId":"2Kc7X1ErDwNQC3mDSzcj2r:2:NEW_SCHEMA_58906353:1.0",
		"type":"CL",
		"tag":"TAG_1",
		"value":{
			"primary":
				{
					"n":"114351986781850739303731927380203402999098269980714037314251113030334656265871709647066791871911223301494528985182940136271565422749649557960004931814103009029733583653015415058729735948968643891217504062569211198683232785559604225075770261432789380861562914440954612182709645333040344110861572682479276273142749932539754452322339633609834159254352611401358466321761471330640054897494238955476899908070070361341619151353045559751402458519854515745183844379712103091543507400422378370355643206447474935832262402237676157940453039766945866619593378491975661408993580741975991704790910675212836114570878937385536741240469","s":"31709631757449918334779137782015080073295053322736405568248206683880709171581786219175755032654289326429594151604974542772479564118714991128114195702524074664728499439192922923447972959261025123487524165136002683835685312396921757867211198611062337367781119789732889820683594639802614428404839292427045186592186376430133296980565403835200091455760550155342562241135705547390154326253551443311233654405275454352469155124427769521657093240947965912272734029661351537725893007773490270205830039503257093200820483663293079052787798381174464715067828789912652332118075432384031240120250649702459279523437925745586018358293",
					"r":{
						"master_secret":"12837988216805632351836502945529578102474557599830967972186207065342685954661648856444421476739667585880031389901024111606799727679248375503534255029537628134340044725081999331175304118567623592026640121065134618663221222181660793523789155609571705183719400775990624021197033214036293996853852857207793762674954939880356614988643856549169115520729781191928393013818919219710385509503729227772128692771910790952034736988280059975119280471359514698093830911966671305913679893961756555349193757551140869441406104544569514488098827228935370871710865416119729343110888522396967560330824755473226246077646585163092579723997","email":"93425314089627999818380533787916809209950740950115820021698943716824532558504838247429084025859644946108373784077910321990397825958470829903294725572816234028190196233941388974321785608037467454430012965788880762183536516406638266343493777424188670573876471142211529285671102219939725377818825746898788375703752623984954453732097535944784566297565616884235473974986796671776425377653568554907233830054733287299779748698319679936358202674100156444583765360920724497237329070042270400326892569611134078842388566552888018013016328624195238117066277964740179687034483497565901714875779198103745484659612191720833693616419"
					},"rctxt":"46100874202549010803287462701390231538257931653160584184367700127275896388140693833771625618896697329876860989596248847885168075827735101683374121816675949155167772610236941660832320228132192916613114017451310845585770227805652692244483696077415158896125061029230966455216903565350085381390999789333348897790885443229989704193891109598610586340077657863453866106011128362359341607159031864090098629788020258146744618532057657525168512208730932379455201671553618168235060463186010956342168679597490845177583194203313273318014748887565074766626101065464049772677912292696033739144461651128229161259333766516373992681879","z":"89161921686183633829410300995305465542210240490443467775323203752739086291058606612122247243462232747738768545470504969895393341806397903979860270213597582400132467523854689715449859935296842541530600458474635876438328705288767879114989512876635995608891378421539351744840475213897790574244733289308957793165323404889405312075964681296404235068158976273816793589717808433803477472546113698906028117688362894051983899751264642818114933575695609005764566927091234057554668558126118496561698874640184178922144296390770694820558052595258564351424575080116398365522802635360189170934064550963860084739798771584807929625024"
				}
		}
	}`

const immuTxnIDForSchema = "2Kc7X1ErDwNQC3mDSzcj2r:2:NEW_SCHEMA_58906353:1.0"
const immuSchemaDataToWrite = `
	{
		"ver":"1.0",
		"id":"2Kc7X1ErDwNQC3mDSzcj2r:2:NEW_SCHEMA_58906353:1.0",
		"name":"NEW_SCHEMA_58906353",
		"version":"1.0",
		"attrNames":["email"],
		"seqNo":null
	}`

const immuTxnIDForNym = "2TEvwu4PeDbfzvAy6e3HgQ"
const immuNymDataToWrite = "UFg64enRj3o3w8arQAJn9U"

var immuURL string
var immuPort string
var userName string
var password string

var storedKey []byte
var storedValue []byte

// create immuClient mock
type mockImmuClient struct {
	immuclient.ImmuClient
}

// Override the real immuclient.Set() function. Can be used to return also errors if needed
func (m *mockImmuClient) Set(ctx context.Context, key []byte, value []byte) (*schema.TxMetadata, error) {
	// store values
	storedKey = key
	storedValue = value
	// Set test data to return. This is how the real data looks like
	var txData schema.TxMetadata
	txData.Id = 108
	txData.PrevAlh = []byte("E+\x1e\x85\x85 X\x1d\x87\x8a\x03\xb1\xf2\xb1\xf5\x9eh\xa2\xf2_5{1Ӎ\x03Bٵڳ\xd9")
	txData.Ts = 1614767958
	txData.EH = []byte("BA\xaab\x9a{Y\xa4\xad\xd9\xee\xa4fn^^Q\x14d\x87k4%\xdcލC\xd6Ԁ\xc7(")
	txData.BlTxId = 107
	txData.BlRoot = []byte("q\xb7(<U]\xba\xad\x8b\xf1\x1cB\x83E\xe6`\xf9\xc3\x12\xe9y\x05\xf9+[\xfawS\xab\xa0\x92I")
	return &txData, nil
}

// Override the real immuclient.Get() function. Can be used to return also errors if needed
func (m *mockImmuClient) Get(ctx context.Context, key []byte) (*schema.Entry, error) {
	// Set test data to return. This is how the real data looks like
	var entryData schema.Entry
	entryData.Tx = 117
	entryData.Key = storedKey
	entryData.Value = storedValue
	return &entryData, nil
}

var testImmuClient = &mockImmuClient{}

// Store the current env setting before running tests
func TestImmuLedger_StartMockedImmuledger(t *testing.T) {

	// Read ImmuDB related credential, Url and port data from env
	// As these are now stored in variables, the env variables can
	// be manipulated if needed for testing against ImmuDB running anywhere
	// Not used for anything at the moment though
	immuURL = os.Getenv("ImmuUrl")
	immuPort = os.Getenv("ImmuPort")
	userName = os.Getenv("ImmuUsrName")
	password = os.Getenv("ImmuPasswd")
}
func TestImmuLedger_Open(t *testing.T) {
	ok := immuLedger.Open("FINDY_IMMUDB_LEDGER")
	assert.True(t, ok)
}

// Test CRED DEF writing and reading
func TestImmuLedger_CRedDef(t *testing.T) {
	ok := immuLedger.Open("FINDY_IMMUDB_LEDGER")
	if !testAgainstRealImmuDB {
		immuLedger.MockImmuClientForTesting(testImmuClient)
	}
	assert.True(t, ok)
	err := immuLedger.Write(immuTxnIDForClaim, immuClaimDataToWrite)
	assert.NoError(t, err)
	// clear MemCache to test reading from ImmuDB / mocked Immu
	immuLedger.ResetMemCache()
	name, value, err := immuLedger.Read(immuTxnIDForClaim)
	assert.NoError(t, err)
	assert.Equal(t, immuTxnIDForClaim, name)
	assert.Equal(t, CleanDataString(immuClaimDataToWrite), value)

	// Read from mem cache
	for i := 0; i < 100; i++ {
		name, value, err := immuLedger.Read(immuTxnIDForClaim)
		assert.NoError(t, err)
		assert.Equal(t, immuTxnIDForClaim, name)
		assert.Equal(t, CleanDataString(immuClaimDataToWrite), value)
	}
	immuLedger.Close()
}

// Test SCHEMA writing and reading
func TestImmuLedger_Schema(t *testing.T) {
	ok := immuLedger.Open("FINDY_IMMUDB_LEDGER")
	if !testAgainstRealImmuDB {
		immuLedger.MockImmuClientForTesting(testImmuClient)
	}
	assert.True(t, ok)
	err := immuLedger.Write(immuTxnIDForSchema, immuSchemaDataToWrite)
	assert.NoError(t, err)
	// clear MemCache to test reading from ImmuDB / mocked Immu
	immuLedger.ResetMemCache()
	name, value, err := immuLedger.Read(immuTxnIDForSchema)
	assert.NoError(t, err)
	assert.Equal(t, immuTxnIDForSchema, name)
	assert.Equal(t, CleanDataString(immuSchemaDataToWrite), value)

	// Read from mem cache
	for i := 0; i < 100; i++ {
		name, value, err := immuLedger.Read(immuTxnIDForSchema)
		assert.NoError(t, err)
		assert.Equal(t, immuTxnIDForSchema, name)
		assert.Equal(t, CleanDataString(immuSchemaDataToWrite), value)
	}
	immuLedger.Close()
}

// Test NYM writing and reading
func TestImmuLedger_Nym(t *testing.T) {
	ok := immuLedger.Open("FINDY_IMMUDB_LEDGER")
	if !testAgainstRealImmuDB {
		immuLedger.MockImmuClientForTesting(testImmuClient)
	}
	assert.True(t, ok)
	err := immuLedger.Write(immuTxnIDForNym, immuNymDataToWrite)
	assert.NoError(t, err)
	// clear MemCache to test reading from ImmuDB / mocked Immu
	immuLedger.ResetMemCache()
	name, value, err := immuLedger.Read(immuTxnIDForNym)
	assert.NoError(t, err)
	assert.Equal(t, immuTxnIDForNym, name)
	assert.Equal(t, CleanDataString(immuNymDataToWrite), value)

	// Read from mem cache
	for i := 0; i < 100; i++ {
		name, value, err := immuLedger.Read(immuTxnIDForNym)
		assert.NoError(t, err)
		assert.Equal(t, immuTxnIDForNym, name)
		assert.Equal(t, CleanDataString(immuNymDataToWrite), value)
	}
	immuLedger.Close()
}

// Restores the env settings
func TestImmuLedger_RestoreEnv(t *testing.T) {
	// restore env
	os.Setenv("ImmuUrl", immuURL)
	os.Setenv("ImmuPort", immuPort)
	os.Setenv("ImmuUsrName", userName)
	os.Setenv("ImmuPasswd", password)
}
