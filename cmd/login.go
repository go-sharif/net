/*
Copyright © 2024 Saman Dehestani <github.com/drippypale>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sharif/net/internal/http"
	cui "github.com/go-sharif/net/internal/ui"
	"github.com/go-sharif/net/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(loginCmd)

	// move to the rootCmd
	loginCmd.Flags().StringP("username", "u", "", "Username to login")
	loginCmd.Flags().StringP("password", "p", "", "Password to login")
	loginCmd.Flags().BoolP("alive", "a", false, "Stay alive to monitor session")

	// remove later to use config instead (using these flags are optional then)
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the network",
	Run: func(cmd *cobra.Command, args []string) {

		// First load the config values. They will be overrided it specified by the user.
		alive := viper.GetBool("alive")

		// Set the passed Flags and Args
		useIP, _ := cmd.Flags().GetBool("use-ip")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		alive, _ = cmd.Flags().GetBool("alive")

		lh := &http.LoginHandler{
			Username: username,
			Password: password,
			UseIP:    useIP,
		}
		if err := lh.Init(); err != nil {
			log.Fatalln(err)
		}

		// Start the
		if loginStatus, err := lh.Login(); loginStatus != 200 || err != nil {
			log.Fatalf("Failed to login: %s", err)
		}
		log.Println("Successfully logined ...")

		if alive {
			// Checking the internet connection and relogin, needs the
			// super user permissions for now.
			isRoot := util.IsRoot()

			// Channels and Contexts
			uiQuitChan := make(chan struct{})
			pingChan := make(chan string)
			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()

			sh := &http.SessionStatusHandler{
				UseIP: useIP,
			}
			// UI init
			luih := &cui.LoginHandler{}
			luih.Init()
			defer luih.Close()

			luih.AddLog("Successfully logined ...", cui.LogSucc)
			if !isRoot {
				luih.AddPing("Checking the internet connection needs the sudo permissions ...", cui.LogErr)
			}

			if isRoot {
				go http.PingTicker(ctx, cancelFunc, pingChan, "8.8.8.8")
				go reloginHandler(sh, lh, uiQuitChan, luih, ctx)
			}
			go luih.PollEvents(uiQuitChan)
			go statusTicker(sh, uiQuitChan, pingChan, luih, ctx)

			// Wait for the proper signals to exit
			// TODO: Replace this with the Context
			<-uiQuitChan
		}
	},
}

func reloginHandler(
	sh *http.SessionStatusHandler,
	lh *http.LoginHandler,
	uiQuitChan chan struct{},
	luih *cui.LoginHandler,
	ctx context.Context,
) {
	for {
		select {
		case <-uiQuitChan:
			return
		case <-ctx.Done():
			for i := 3; i > 0; i-- {
				luih.AddLog(fmt.Sprintf("Internet connection lost ... Retry in %v ...", i), cui.LogErr)
				time.Sleep(1 * time.Second)
			}

			if loginStatus, err := lh.Login(); loginStatus != 200 || err != nil {
				luih.AddLog(fmt.Sprintf("Failed to login: %s", err), cui.LogErr)
				time.Sleep(1 * time.Second)
				os.Exit(1)
			} else {
				luih.AddLog("Successfully logined  again ...", cui.LogSucc)
			}

			// renew channels
			uiQuitChan = make(chan struct{})
			pingChan := make(chan string)

			ctxN, cancelFuncN := context.WithCancel(context.Background())
			ctx = ctxN

			go statusTicker(sh, uiQuitChan, pingChan, luih, ctx)
			go http.PingTicker(ctx, cancelFuncN, pingChan, "8.8.8.8")

			time.Sleep(500 * time.Millisecond)
		}
	}

}

func statusTicker(
	sh *http.SessionStatusHandler,
	uiQuitChan chan struct{},
	pingChan chan string,
	luih *cui.LoginHandler,
	ctx context.Context,
) {
	if err := sh.Init(); err != nil {
		luih.AddLog(err.Error(), cui.LogErr)
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-uiQuitChan:
			return
		case <-ctx.Done():
			return
		case p := <-pingChan:
			luih.AddPing(p, cui.LogInfo)
		case <-ticker.C:
			if statusCode, ss, err := sh.GetSessionStatus(true); err == nil && statusCode == 200 && ss != nil {
				luih.UpdateStatusTable(ss)
				if diff := sh.Diff(); diff != nil {
					luih.AddBytesData(float64(diff.BytesUp), float64(diff.BytesDown))
				}
				luih.Refresh()
			}
		}
	}
}
