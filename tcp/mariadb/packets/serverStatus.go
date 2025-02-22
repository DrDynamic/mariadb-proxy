package packets

const SERVER_STATUS_IN_TRANS = 1                  // A transaction is currently active
const SERVER_STATUS_AUTOCOMMIT = 2                // Autocommit mode is set
const SERVER_MORE_RESULTS_EXISTS = 8              // More results exists (more packets will follow)
const SERVER_QUERY_NO_GOOD_INDEX_USED = 16        // Set if EXPLAIN would've shown Range checked for each record
const SERVER_QUERY_NO_INDEX_USED = 32             // The query did not use an index
const SERVER_STATUS_CURSOR_EXISTS = 64            // When using COM_STMT_FETCH, indicate that current cursor still has result
const SERVER_STATUS_LAST_ROW_SENT = 128           // When using COM_STMT_FETCH, indicate that current cursor has finished to send results
const SERVER_STATUS_DB_DROPPED = 1 << 8           // Database has been dropped
const SERVER_STATUS_NO_BACKSLASH_ESCAPES = 1 << 9 // Current escape mode is "no backslash escape"
const SERVER_STATUS_METADATA_CHANGED = 1 << 10    // A DDL change did have an impact on an existing PREPARE (an automatic reprepare has been executed)
const SERVER_QUERY_WAS_SLOW = 1 << 11             // The query was slower than long_query_time
const SERVER_PS_OUT_PARAMS = 1 << 12              // This resultset contain stored procedure output parameter
const SERVER_STATUS_IN_TRANS_READONLY = 1 << 13   // Current transaction is a read-only transaction
const SERVER_SESSION_STATE_CHANGED = 1 << 14      // Session state change. See Session change type for more information
