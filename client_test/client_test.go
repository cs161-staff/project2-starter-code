package client_test

// You MUST NOT change these default imports.  ANY additional imports it will
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports. Normally, you will want to avoid underscore imports
	// unless you know exactly what you are doing. You can read more about
	// underscore imports here: https://golangdocs.com/blank-identifier-in-golang
	_ "encoding/hex"
	_ "encoding/json"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect(). You can read more
	// about dot imports here:
	// https://stackoverflow.com/questions/6478962/what-does-the-dot-or-period-in-a-go-import-statement-do
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	// The client implementation is intentionally defined in a different package.
	// This forces us to follow best practice and write tests that only rely on
	// client API that is exported from the client package, and avoid relying on
	// implementation details private to the client package.
	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	// We are using 2 libraries to help us write readable and maintainable tests:
	//
	// (1) Ginkgo, a Behavior Driven Development (BDD) testing framework that
	//             makes it easy to write expressive specs that describe the
	//             behavior of your code in an organized manner; and
	//
	// (2) Gomega, an assertion/matcher library that allows us to write individual
	//             assertion statements in tests that read more like natural
	//             language. For example "Expect(ACTUAL).To(Equal(EXPECTED))".
	//
	// In the Ginko framework, a test case signals failure by calling Ginkgo’s
	// Fail(description string) function. However, we are using the Gomega library
	// to execute our assertion statements. When a Gomega assertion fails, Gomega
	// calls a GomegaFailHandler, which is a function that must be provided using
	// gomega.RegisterFailHandler(). Here, we pass Ginko's Fail() function to
	// Gomega so that Gomega can report failed assertions to the Ginko test
	// framework, which can take the appropriate action when a test fails.
	//
	// This is the sole connection point between Ginkgo and Gomega.
	RegisterFailHandler(Fail)

	RunSpecs(t, "Client Tests")
}

// ================================================
// Here are some optional global variables that can be used throughout the test
// suite to make the tests more readable and maintainable than defining these
// values in each test. You can add more variables here if you want and think
// they will help keep your code clean!
// ================================================
const someFilename = "file1.txt"
const someOtherFilename = "file2.txt"
const nonExistentFilename = "thisFileDoesNotExist.txt"

const aliceUsername = "Alice"
const alicePassword = "AlicePassword"
const bobUsername = "Bob"
const bobPassword = "BobPassword"
const nilufarUsername = "Nilufar"
const nilufarPassword = "NilufarPassword"
const olgaUsername = "Olga"
const olgaPassword = "OlgaPassword"
const marcoUsername = "Marco"
const marcoPassword = "MarcoPassword"

const nonExistentUsername = "NonExistentUser"

var alice *client.User
var bob *client.User
var nilufar *client.User
var olga *client.User
var marco *client.User

var someFileContent []byte
var someShortFileContent []byte
var someLongFileContent []byte

// ================================================
// The top level Describe() contains all tests in
// this test suite in nested Describe() blocks.
// ================================================

var _ = Describe("Client Tests", func() {
	BeforeEach(func() {
		// This top-level BeforeEach will be run before each test.
		//
		// Resets the state of Datastore and Keystore so that tests do not
		// interfere with each other.
		userlib.DatastoreClear()
		userlib.KeystoreClear()

		// You can set this to false!
		userlib.SetDebugStatus(true)
	})

	BeforeEach(func() {
		// This top-level BeforeEach will be run before each test.
		//
		// Byte slices cannot be constant, so this BeforeEach resets the content of
		// each global variable to a predefined value, which allows tests to rely on
		// the expected value of these variables.
		someShortFileContent = []byte("some short file content")
		someFileContent = someShortFileContent
		someLongFileContent = []byte("some LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOONG file content")
	})

	Describe("Creating users", func() {
		It("should not error when creating a new user", func() {
			_, err := client.InitUser("Alice", "password")
			Expect(err).To(BeNil(), "Failed to initialized user Alice.")
		})

		It("should error if a username is already taken by another user", func() {
			// TODO: implement this function
		})

		It("should error if a user does not exist with that username", func() {
			// TODO: implement function
		})

		// TODO: you probably want more test cases about creating users here
	})

	Describe("Single user storage", func() {
		var alice *client.User

		BeforeEach(func() {
			// This BeforeEach will run before each test in this Describe block.
			alice, _ = client.InitUser("Alice", "some password")
		})

		It("should upload content without erroring", func() {
			content := []byte("This is a test")
			err := alice.StoreFile("file1", content)
			Expect(err).To(BeNil(), "Failed to upload content to a file", err)
		})

		It("should download the expected content that was previously uploaded", func() {
			uploadedContent := []byte("This is a test")
			alice.StoreFile(someFilename, uploadedContent)
			downloadedContent, _ := alice.LoadFile(someFilename)
			Expect(downloadedContent).To(BeEquivalentTo(uploadedContent),
				"Downloaded content is not the same as uploaded content",
				downloadedContent,
				uploadedContent)
		})

		It("should error when trying to download a file that does not exist", func() {
			_, err := alice.LoadFile(nonExistentFilename)
			Expect(err).ToNot(BeNil(), "Was able to load a non-existent file without error.")
		})

		// TODO: you probably want more test cases for store/load/append with a
		// 			 single user here
	})

	Describe("Sharing files", func() {

		BeforeEach(func() {
			// Initialize each user to ensure the variable has the expected value for
			// the tests in this Describe() block.
			alice, _ = client.InitUser(aliceUsername, alicePassword)
			bob, _ = client.InitUser(bobUsername, bobPassword)
			nilufar, _ = client.InitUser(nilufarUsername, nilufarPassword)
			olga, _ = client.InitUser(olgaUsername, olgaPassword)
			marco, _ = client.InitUser(marcoUsername, marcoPassword)
		})

		It("should share a file without erroring", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")

			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Bob downloaded was not the same as what Alice uploaded.")
		})

		// TODO: you probably want more test cases for sharing files here
	})

	// TODO: you probably want more Describe() blocks to contain tests related to
	//       logical test groupings other than the ones suggested above
})
