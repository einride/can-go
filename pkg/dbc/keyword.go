package dbc

// Keyword represents a DBC keyword.
type Keyword string

const (
	KeywordAttribute               Keyword = "BA_DEF_"
	KeywordAttributeDefault        Keyword = "BA_DEF_DEF_"
	KeywordAttributeValue          Keyword = "BA_"
	KeywordBitTiming               Keyword = "BS_"
	KeywordComment                 Keyword = "CM_"
	KeywordEnvironmentVariable     Keyword = "EV_"
	KeywordEnvironmentVariableData Keyword = "ENVVAR_DATA_"
	KeywordMessage                 Keyword = "BO_"
	KeywordMessageTransmitters     Keyword = "BO_TX_BU_"
	KeywordNewSymbols              Keyword = "NS_"
	KeywordNodes                   Keyword = "BU_"
	KeywordSignal                  Keyword = "SG_"
	KeywordSignalMultiplexValue    Keyword = "SG_MUL_VAL_"
	KeywordSignalGroup             Keyword = "SIG_GROUP_"
	KeywordSignalType              Keyword = "SGTYPE_"
	KeywordSignalValueType         Keyword = "SIG_VALTYPE_"
	KeywordValueDescriptions       Keyword = "VAL_"
	KeywordValueTable              Keyword = "VAL_TABLE_"
	KeywordVersion                 Keyword = "VERSION"
)
