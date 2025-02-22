package packets

const CLIENT_MYSQL = 1 // Set by older MariaDB versions. MariaDB 10.2 leaves this bit unset to permit MariaDB identification and indicate support for extended capabilities. (MySQL named this CLIENT_LONG_PASSWORD)
const FOUND_ROWS = 2
const CONNECT_WITH_DB = 8         // One can specify db on connect
const COMPRESS = 32               // Can use compression protocol
const LOCAL_FILES = 128           // Can use LOAD DATA LOCAL
const IGNORE_SPACE = 256          //	Ignore spaces before '('
const CLIENT_PROTOCOL_41 = 1 << 9 // 4.1 protocol
const CLIENT_INTERACTIVE = 1 << 10
const SSL = 1 << 11 // Can use SSL
const TRANSACTIONS = 1 << 13
const SECURE_CONNECTION = 1 << 15                   // 4.1 authentication
const MULTI_STATEMENTS = 1 << 16                    // Enable/disable multi-stmt support
const MULTI_RESULTS = 1 << 17                       // Enable/disable multi-results
const PS_MULTI_RESULTS = 1 << 18                    // Enable/disable multi-results for PrepareStatement
const PLUGIN_AUTH = 1 << 19                         // Client supports plugin authentication
const CONNECT_ATTRS = 1 << 20                       // Client send connection attributes
const PLUGIN_AUTH_LENENC_CLIENT_DATA = 1 << 21      // Enable authentication response packet to be larger than 255 bytes
const CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS = 1 << 22 // Client can handle expired passwords
const CLIENT_SESSION_TRACK = 1 << 23                // Enable/disable session tracking in OK_Packet
const CLIENT_DEPRECATE_EOF = 1 << 24                // EOF_Packet deprecation :
// * OK_Packet replace EOF_Packet in end of Resulset when in text format
// * EOF_Packet between columns definition and resultsetRows is deleted
const CLIENT_OPTIONAL_RESULTSET_METADATA = 1 << 25 // Not use for MariaDB
const CLIENT_ZSTD_COMPRESSION_ALGORITHM = 1 << 26  // Support zstd protocol compression
const CLIENT_CAPABILITY_EXTENSION = 1 << 29        // Reserved for future use. (Was CLIENT_PROGRESS Client support progress indicator before 10.2)
const CLIENT_SSL_VERIFY_SERVER_CERT = 1 << 30      // Client verify server certificate. deprecated, client have options to indicate if server certifiate must be verified
const CLIENT_REMEMBER_OPTIONS = 1 << 31
const MARIADB_CLIENT_PROGRESS = 1 << 32             // Client support progress indicator (since 10.2)
const MARIADB_CLIENT_COM_MULTI = 1 << 33            // Permit COM_MULTI protocol
const MARIADB_CLIENT_STMT_BULK_OPERATIONS = 1 << 34 // Permit bulk insert
const MARIADB_CLIENT_EXTENDED_METADATA = 1 << 35    // Add extended metadata information
const MARIADB_CLIENT_CACHE_METADATA = 1 << 36       // Permit skipping metadata
const MARIADB_CLIENT_BULK_UNIT_RESULTS = 1 << 37    // when enable, indicate that Bulk command can use STMT_BULK_FLAG_SEND_UNIT_RESULTS flag that permit to return a result-set of all affected rows and auto-increment values
