/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Http Dumper",
	Long:  `Run Http Dumper`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if destUrl == "" {
			return fmt.Errorf("desination host url ")
		}
		fmt.Println("run called")
		target, err := url.Parse(destUrl)
		if err != nil {
			return fmt.Errorf("unable to parse destination host url: %w", err)
		}
		log.Printf("forwarding to -> %s%s\n", target.Scheme, target.Host)

		proxy := httputil.NewSingleHostReverseProxy(target)

		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			// https://stackoverflow.com/questions/38016477/reverse-proxy-does-not-work
			// https://forum.golangbridge.org/t/explain-how-reverse-proxy-work/6492/7
			// https://stackoverflow.com/questions/34745654/golang-reverseproxy-with-apache2-sni-hostname-error

			requestID, _ := uuid.NewV4()
			req.Host = req.URL.Host
			wr := &wrappedResponseWriter{rw: w}
			bR, _ := httputil.DumpRequest(req, true)
			proxy.ServeHTTP(wr, req)
			log.Printf("============= %s =============\nRequest >>>>>>>>>>>>>>\n%s\nRsponse <<<<<<<<<<<<\n%s", requestID.String(), string(bR), wr.dump())
		})

		err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			panic(err)
		}
		return nil
	},
}

type wrappedResponseWriter struct {
	rw         http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func (w *wrappedResponseWriter) Header() http.Header {
	return w.rw.Header()
}
func (w *wrappedResponseWriter) Write(buf []byte) (int, error) {
	w.buf.Write(buf)
	return w.rw.Write(buf)
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.rw.WriteHeader(statusCode)
}

func (w *wrappedResponseWriter) dump() string {
	headers := make([]string, 0, len(w.rw.Header()))
	for k, v := range w.rw.Header() {
		headers = append(headers, fmt.Sprintf("%s: %s", k, v))
	}
	return fmt.Sprintf("Status: %d\n%s\n\n%s", w.statusCode, strings.Join(headers, "\n"), string(w.buf.Bytes()))
}

var (
	destUrl string
	port    uint64
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().Uint64Var(&port, "port", 8397, "A listening port")
	runCmd.PersistentFlags().StringVar(&destUrl, "dest-url", "", "A destination host url")
}
