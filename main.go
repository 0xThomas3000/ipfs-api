package main

import (
	"fmt"
	"time"

	"io"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

// Paste here the local path of your computer where the file will be downloaded
const YourLocalPath = "~/Downloads"

// Paste here your public key of IPFS node
const YourPublicKey = ""

// Create a new file with a text as param
func addFile(sh *shell.Shell, text string) (string, error) {
	// + Returns a CID of the file if added correctly to IPFS
	// + Have to provide the Reader to the Add method (can create a reader from an external file or a string)
	return sh.Add(strings.NewReader(text))
}

func readFile(sh *shell.Shell, cid string) (*string, error) {
	// + use the Cat() of the shell
	reader, err := sh.Cat(fmt.Sprintf("/ipfs/%s", cid))
	if err != nil { // if the error is populated, we stop reading the file, return an empty string
		return nil, fmt.Errorf("error reading the file: %s", err.Error())
	}

	// The "reader" is an abstraction for the content of a file, so we need go from a reader to a string
	// 1. Take the bytes from the reader
	// 2. Check for error
	// 3. Convert the bytes to the string
	// 4. Return the text variable
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading the bytes: %s", err.Error()) // stop the flow and return an empty string
	}

	text := string(bytes) //

	return &text, nil
}

func downloadFile(sh *shell.Shell, cid string) error {
	return sh.Get(cid, YourLocalPath)
}

// Publish an IPNS record
func addToIPNS(sh *shell.Shell, cid string) error {

	// Use publish with details method from the client which expects 5 params
	// + lifetime: the amount of time that the record will be valid for
	// + ttl: the amount of time that IPFS will cache the record
	// + last param: indicates if the client should check if the IPNS record already exists
	var lifetime time.Duration = 50 * time.Hour  // 50 hours
	var ttl time.Duration = 1 * time.Microsecond // 1 microsecond so that IPNS doesn't cache the record

	// Setting the low cache time: if we update the record, it'll be replaced immediately
	_, err := sh.PublishWithDetails(cid, YourPublicKey, lifetime, ttl, true)
	return err
}

// To retrieve the IPFS path which is stored in an IPNS record
func resolveIPNS(sh *shell.Shell) (string, error) {
	return sh.Resolve(YourPublicKey)
}

/* main(): calls the helper functions and handles results. */
func main() {
	// 1. The first thing: need to create a new client.
	// We provide the URL of IPFS API as a param with a default port
	sh := shell.NewShell("localhost:5001") // Returns a client that contains all methods to interact with IPFS node

	err := performChecks(sh)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	/* => We had a client -> pass it to the different helper functions to Create, Read, Download an IPFS file */
	// 2. Add the "Hello from Launchpad!" text to IPFS
	fmt.Println("Adding file to IPFS")
	cid, err := addFile(sh, "Hello from Launchpad!")
	if err != nil {
		fmt.Println("Error adding file to IPFS:", err.Error())
		return
	}
	fmt.Println("File added with CID:", cid)

	separator()

	// 3. Read the file by using the generated CID
	fmt.Println("Reading file")
	text, err := readFile(sh, cid)
	if err != nil {
		fmt.Println("Error reading the file:", err.Error())
		return
	}
	fmt.Println("Content of the file:", *text)

	separator()

	// 4. Download the file to your computer
	fmt.Println("Downloading file")
	err = downloadFile(sh, cid)
	if err != nil {
		fmt.Println("Error downloading file:", err.Error())
		return
	}
	fmt.Println("File donwloaded")

	separator()

	// 5. Publish the file to IPNS
	fmt.Println("Adding file to IPNS")
	err = addToIPNS(sh, cid)
	if err != nil {
		fmt.Println("Error publishing to IPNS:", err.Error())
		return
	}
	fmt.Println("File added to IPNS")

	separator()

	// 6. Resolve IPNS based on your public key
	fmt.Println("Resolving file in IPNS")
	result, err := resolveIPNS(sh)
	if err != nil {
		fmt.Println("Error resolving IPNS:", err.Error())
		return
	}

	fmt.Println("IPNS is pointing to:", result)
}
