package main

import "os"

// SQLConnectString localize and protect connection string
var SQLConnectString = os.Getenv("DBMYSQL")
