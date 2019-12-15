package cli

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/posmint/types"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"strconv"
	"strings"
	"syscall"
)

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(createCmd)
	accountsCmd.AddCommand(deleteCmd)
	accountsCmd.AddCommand(listCmd)
	accountsCmd.AddCommand(showCmd)
	accountsCmd.AddCommand(updatePassphraseCmd)
	accountsCmd.AddCommand(signCmd)
	accountsCmd.AddCommand(importArmoredCmd)
	accountsCmd.AddCommand(importCmd)
	accountsCmd.AddCommand(exportCmd)
	accountsCmd.AddCommand(exportRawCmd)
	accountsCmd.AddCommand(sendTxCmd)
	accountsCmd.AddCommand(sendTxCmd)
}

// accountsCmd represents the accounts namespace command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "The accounts namespace",
	Long: `The accounts namespace handles all account related interactions, 
from creating and deleting accounts, to importing and exporting accounts.`,
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new account",
	Long:  `Creates and persists a new account in the Keybase. Will prompt the user for a passphrase to encrypt the generated keypair.`,
	Run: func(cmd *cobra.Command, args []string) {
		kb, err := app.GetKeybase()
		if err != nil {
			panic(err)
		}
		fmt.Print("Enter Password: ")
		kp, err := kb.Create(credentials())
		if err != nil {
			panic(err)
		}
		fmt.Printf("Account generated succesfully:\nPublic Key: %s", hex.EncodeToString(kp.PubKey.Bytes()))
	},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete  <address>",
	Short: "Delete an account",
	Long:  `Deletes a keypair from the keybase. Will prompt the user for the account passphrase`,
	Run: func(cmd *cobra.Command, args []string) {
		kb, err := app.GetKeybase()
		if err != nil {
			panic(err)
		}
		addr, err := types.AccAddressFromHex(args[0])
		if err != nil {
			panic(err)
		}
		fmt.Print("Enter Password: ")
		err = kb.Delete(addr, credentials())
		if err != nil {
			panic(err)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long: `Lists all the account addresses stored in the keybase. Example output:
(0) 0xb3746D30F2A579a2efe7F2F6E8E06277a78054C1
(1) 0xab514F27e98DE7E3ecE3789b511dA955C3F09Bbc`,
	Run: func(cmd *cobra.Command, args []string) {
		kb, err := app.GetKeybase()
		if err != nil {
			panic(NewBeforeInitError(errors.New("nil keybase")))
		}
		kp, err := kb.List()
		if err != nil {
			panic(err)
		}
		for i, key := range kp {
			fmt.Printf("(%d) %s", i, key)
		}
	},
}

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show  <address>",
	Short: "Shows a pubkey for address",
	Long:  `Lists an account address and public key`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kb, err := app.GetKeybase()
		if err != nil {
			panic(NewBeforeInitError(errors.New("nil keybase")))
		}
		addr, err := types.AccAddressFromHex(args[0])
		if err != nil {
			panic(err)
		}
		kp, err := kb.Get(addr)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Address: %s\nPublic Key: %s\n", kp.GetAddress().String(), hex.EncodeToString(kp.PubKey.Bytes()))
	},
}

// updatePassphraseCmd represents the updatePassphrase command
var updatePassphraseCmd = &cobra.Command{
	Use:   "updatePassphrase <address>",
	Short: "Update account passphrase",
	Long:  `Updates the passphrase for the indicated account. Will prompt the user for the current account passphrase and the new account passphrase.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kb, err := app.GetKeybase()
		if err != nil {
			panic(err)
		}
		addr, err := types.AccAddressFromHex(args[0])
		if err != nil {
			panic(err)
		}
		fmt.Print("Enter Password: ")
		oldpass := credentials()
		fmt.Print("Enter New Password: ")
		newpass := credentials()
		err = kb.Update(addr, oldpass, newpass)
		if err != nil {
			panic(err)
		}
		fmt.Println("Successfully updated account: " + addr.String())
	},
}

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "pocket accounts sign <address> <msg>",
	Short: "Sign a message with an account",
	Long:  `Digitally signs the specified <msg> using the specified <address> account credentials. Will prompt the user for the account passphrase.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		kb, err := app.GetKeybase()
		if err != nil {
			panic(err)
		}
		addr, err := types.AccAddressFromHex(args[0])
		if err != nil {
			panic(err)
		}
		msg, err := hex.DecodeString(args[1])
		if err != nil {
			panic(err)
		}
		sig, _, err := kb.Sign(addr, credentials(), msg)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Original Message: %s\n Signature: %s", args[1], sig)
	},
}

var importArmoredCmd = &cobra.Command{
	Use:   "import-armored <armor>",
	Short: "Import keypair using armor",
	Long:  `Imports an account using the Encrypted ASCII armored <armor> string. Will prompt the user for a decryption passphrase of the <armor> string and for an encryption passphrase to store in the Keybase.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kb, err := app.GetKeybase()
		if err != nil {
			panic(err)
		}
		fmt.Println("Enter decrypt pass")
		dPass := credentials()
		fmt.Println("Enter encrypt pass")
		ePass := credentials()
		kp, err := kb.ImportPrivKey(args[0], dPass, ePass)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Account imported succesfully:\nPublic Key: %s", hex.EncodeToString(kp.PubKey.Bytes()))
	},
}

var exportCmd = &cobra.Command{
	Use:   "export <address>",
	Short: "Export an account",
	Args:  cobra.ExactArgs(1),
	Long: `Exports the account with <address>, encrypted and ASCII armored. 
Will prompt the user for the account passphrase and for an encryption passphrase for the exported account.
`,
	Run: func(cmd *cobra.Command, args []string) {
		kb, err := app.GetKeybase()
		if err != nil {
			panic(err)
		}
		addr, err := types.AccAddressFromHex(args[0])
		if err != nil {
			panic(err)
		}
		fmt.Println("Enter Decrypt Passphrase")
		dPass := credentials()
		fmt.Println("Enter Encrypt Passphrase")
		ePass := credentials()
		pk, err := kb.ExportPrivKeyEncryptedArmor(addr, dPass, ePass)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Exported account: %s", pk)
	},
}

// exportRawCmd represents the exportRaw command
var exportRawCmd = &cobra.Command{
	Use:   "export-raw <address>",
	Short: "Export Plaintext Privkey",
	Args:  cobra.ExactArgs(1),
	Long: `Exports the raw private key in hex format. Will prompt the user for the account passphrase. 
NOTE: THIS METHOD IS NOT RECOMMENDED FOR SECURITY REASONS, USE AT YOUR OWN RISK.*`,
	Run: func(cmd *cobra.Command, args []string) {
		kb, err := app.GetKeybase()
		if err != nil {
			panic(err)
		}
		addr, err := types.AccAddressFromHex(args[0])
		if err != nil {
			panic(err)
		}
		fmt.Println("Enter Decrypt Passphrase")
		dPass := credentials()
		pk, err := kb.ExportPrivateKeyObject(addr, dPass)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Exported account: %s", hex.EncodeToString(pk.Bytes()))
	},
}

// sendTxCmd represents the sendTx command
var sendTxCmd = &cobra.Command{
	Use:   "send-tx <fromAddr> <toAddr> <amount>",
	Short: "Send POKT",
	Long:  `Sends <amount> POKT <fromAddr> to <toAddr>. Prompts the user for <fromAddr> account passphrase.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		amount, err := strconv.Atoi(args[2])
		if err != nil {
			panic(err)
		}
		res, err := app.SendTransaction(args[0], args[1], credentials(), types.NewInt(int64(amount)))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Transaction Submitted: %s", res.TxHash)
	},
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import-raw",
	Short: "import-raw <private-key-hex>",
	Args:  cobra.ExactArgs(1),
	Long: `Imports an account using the provided <private-key-hex>

Will prompt the user for a passphrase to encrypt the generated keypair.
`,
	Run: func(cmd *cobra.Command, args []string) {
		pkBytes, err := hex.DecodeString(args[0])
		kb, err := app.GetKeybase()
		if err != nil {
			panic(err)
		}
		fmt.Println("Enter Encrypt Passphrase")
		ePass := credentials()
		var pk [64]byte
		copy(pk[:], pkBytes)
		kp, err := kb.ImportPrivateKeyObject(pk, ePass)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Account imported succesfully:\nPublic Key: %s", hex.EncodeToString(kp.PubKey.Bytes()))
	},
}

func credentials() string {
	bytePassword, err := terminal.ReadPassword(syscall.Stdin)
	if err == nil {
		fmt.Println("\nPassword typed: " + string(bytePassword))
	}
	return strings.TrimSpace(string(bytePassword))
}