/*
Copyright 2019 The Kubernetes Authors.

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

package certs

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/grantae/certinfo"
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"k8s.io/ingress-nginx/cmd/plugin/kubectl"
	"k8s.io/ingress-nginx/cmd/plugin/request"
	"k8s.io/ingress-nginx/cmd/plugin/util"
)

// CreateCommand creates and returns this cobra subcommand
func CreateCommand(flags *genericclioptions.ConfigFlags) *cobra.Command {
	var pod, deployment, selector, container *string
	cmd := &cobra.Command{
		Use:   "certs",
		Short: "Output the certificate data stored in an ingress-nginx pod",
		RunE: func(cmd *cobra.Command, args []string) error {
			host, err := cmd.Flags().GetString("host")
			if err != nil {
				return err
			}
			pretty, err := cmd.Flags().GetBool("pretty")
			if err != nil {
				return err
			}

			util.PrintError(certs(flags, *pod, *deployment, *selector, *container, host, pretty))
			return nil
		},
	}

	cmd.Flags().String("host", "", "Get the cert for this hostname")
	if err := cobra.MarkFlagRequired(cmd.Flags(), "host"); err != nil {
		util.PrintError(err)
		os.Exit(1)
	}
	pod = util.AddPodFlag(cmd)
	deployment = util.AddDeploymentFlag(cmd)
	selector = util.AddSelectorFlag(cmd)
	container = util.AddContainerFlag(cmd)
	cmd.Flags().Bool("pretty", false, "Pretty print certificates")
	return cmd
}

func certs(flags *genericclioptions.ConfigFlags, podName string, deployment string, selector string, container string, host string, prettyPrint bool) error {
	command := []string{"/dbg", "certs", "get", host}

	pod, err := request.ChoosePod(flags, podName, deployment, selector)
	if err != nil {
		return err
	}

	out, err := kubectl.PodExecString(flags, &pod, container, command)
	if err != nil {
		return err
	}
	if prettyPrint {
		prettyPrintPemBlocks([]byte(out))
	} else {
		fmt.Print(out)
	}
	return nil
}

func prettyPrintPemBlocks(pemBlocks []byte) {
	var demBlock *pem.Block
	for {
		demBlock, pemBlocks = pem.Decode(pemBlocks)
		if demBlock == nil {
			break
		}
		switch demBlock.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(demBlock.Bytes)
			if err != nil {
				fmt.Println(err)
				break
			}
			info, err := certinfo.CertificateText(cert)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("-----BEGIN CERTIFICATE-----")
			fmt.Print(info)
			fmt.Println("-----END CERTIFICATE-----")
		default:
			fmt.Printf("skip %s\n", demBlock.Type)
		}
	}
}
