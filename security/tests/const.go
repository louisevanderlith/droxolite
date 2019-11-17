package tests

import (
	"crypto/rsa"
)

const testClientID string = "tJVgDxKa"
const testClientSecret string = "JFmXTWvaHO"
const testUserName string = "fake@mango.avo"
const testUserPassword string = "password1"

var testPrivKey *rsa.PrivateKey
var testPubKey *rsa.PublicKey
