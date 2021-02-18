package addons

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const txnIDForClaim = "2Kc7X1ErDwNQC3mDSzcj2r:3:CL:2Kc7X1ErDwNQC3mDSzcj2r:2:NEW_SCHEMA_58906353:1.0:TAG_1"
const claimDataToWrite = `
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

const txnIDForSchema = "2Kc7X1ErDwNQC3mDSzcj2r:2:NEW_SCHEMA_58906353:1.0"
const schemaDataToWrite = `
	{
		"ver":"1.0",
		"id":"2Kc7X1ErDwNQC3mDSzcj2r:2:NEW_SCHEMA_58906353:1.0",
		"name":"NEW_SCHEMA_58906353",
		"version":"1.0",
		"attrNames":["email"],
		"seqNo":null
	}`

const txnIDForNym = "2TEvwu4PeDbfzvAy6e3HgQ"
const nymDataToWrite = "UFg64enRj3o3w8arQAJn9U"
const nymDataToWriteJson = `{"id":"2TEvwu4PeDbfzvAy6e3HgQ","data":"UFg64enRj3o3w8arQAJn9U"}`

var authToken string
var baseAddress string
var mockApiLedger *httptest.Server

// Store the current env setting before running tests
func TestApiLedger_StartMockedLedgerAPI(t *testing.T) {

	// Read API Ledger Auth Token and base address from env
	authToken = os.Getenv("AuthToken")
	baseAddress = os.Getenv("BaseAddress")
}
func TestApiLedger_Open(t *testing.T) {
	ok := apiLedger.Open("FINDY_API_LEDGER")
	assert.True(t, ok)
}

// Test CRED DEF writing and reading
func TestApiLedger_CRedDef(t *testing.T) {

	// Comment this mockApiLedger setting out if you want to test against Ledger API
	// ****************************************************************
	mockApiLedger = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(CleanString(claimDataToWrite)))
		}),
	)
	defer mockApiLedger.Close()
	tempBaseAddress := fmt.Sprint(mockApiLedger.URL, "/")
	// Set env to testServer for unit testing against mocked Ledger API
	os.Setenv("AuthToken", "aff62524-e503-11e9-81b4-2a2ae2dbcce4")
	os.Setenv("BaseAddress", tempBaseAddress)
	// ****************************************************************

	ok := apiLedger.Open("FINDY_API_LEDGER")
	assert.True(t, ok)
	err := apiLedger.Write(txnIDForClaim, claimDataToWrite)
	assert.NoError(t, err)
	name, value, err := apiLedger.Read(txnIDForClaim)
	assert.NoError(t, err)
	assert.Equal(t, txnIDForClaim, name)
	assert.Equal(t, CleanString(claimDataToWrite), value)

	// Read from mem cache
	for i := 0; i < 100; i++ {
		name, value, err := apiLedger.Read(txnIDForClaim)
		assert.NoError(t, err)
		assert.Equal(t, txnIDForClaim, name)
		assert.Equal(t, CleanString(claimDataToWrite), value)
	}
}

// Test SCHEMA writing and reading
func TestApiLedger_Schema(t *testing.T) {
	// Comment this mockApiLedger setting out if you want to test against Ledger API
	// ****************************************************************
	mockApiLedger = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(CleanString(schemaDataToWrite)))
		}),
	)
	defer mockApiLedger.Close()
	tempBaseAddress := fmt.Sprint(mockApiLedger.URL, "/")
	// Set env to testServer for unit testing against mocked Ledger API
	os.Setenv("AuthToken", "aff62524-e503-11e9-81b4-2a2ae2dbcce4")
	os.Setenv("BaseAddress", tempBaseAddress)
	// ****************************************************************
	ok := apiLedger.Open("FINDY_API_LEDGER")
	assert.True(t, ok)
	err := apiLedger.Write(txnIDForSchema, schemaDataToWrite)
	assert.NoError(t, err)
	name, value, err := apiLedger.Read(txnIDForSchema)
	assert.NoError(t, err)
	assert.Equal(t, txnIDForSchema, name)
	assert.Equal(t, CleanString(schemaDataToWrite), value)

	// Read from mem cache
	for i := 0; i < 100; i++ {
		name, value, err := apiLedger.Read(txnIDForSchema)
		assert.NoError(t, err)
		assert.Equal(t, txnIDForSchema, name)
		assert.Equal(t, CleanString(schemaDataToWrite), value)
	}
}

// Test NYM writing and reading
func TestApiLedger_Nym(t *testing.T) {
	// Comment this mockApiLedger setting out if you want to test against Ledger API
	// ****************************************************************
	mockApiLedger = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(nymDataToWriteJson))
		}),
	)
	defer mockApiLedger.Close()
	tempBaseAddress := fmt.Sprint(mockApiLedger.URL, "/")
	// Set env to testServer for unit testing against mocked Ledger API
	os.Setenv("AuthToken", "aff62524-e503-11e9-81b4-2a2ae2dbcce4")
	os.Setenv("BaseAddress", tempBaseAddress)
	// ****************************************************************
	ok := apiLedger.Open("FINDY_API_LEDGER")
	assert.True(t, ok)
	err := apiLedger.Write(txnIDForNym, nymDataToWrite)
	assert.NoError(t, err)
	name, value, err := apiLedger.Read(txnIDForNym)
	assert.NoError(t, err)
	assert.Equal(t, txnIDForNym, name)
	assert.Equal(t, CleanString(nymDataToWrite), value)

	// Read from mem cache
	for i := 0; i < 100; i++ {
		name, value, err := apiLedger.Read(txnIDForNym)
		assert.NoError(t, err)
		assert.Equal(t, txnIDForNym, name)
		assert.Equal(t, CleanString(nymDataToWrite), value)
	}
}

// Restores the env settings
func TestApiLedger_RestoreEnv(t *testing.T) {
	// restore env
	os.Setenv("AuthToken", authToken)
	os.Setenv("BaseAddress", baseAddress)
}
