package cmd

import (
	"context"
	"fmt"

	"github.com/b-harvest/indep_node_alarm_go/client"
	"github.com/b-harvest/indep_node_alarm_go/config"
	tmtype "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "indep [config_path]",
		Args:  cobra.ExactArgs(1),
		Short: "Example: $indep ./config.toml",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			indep_config := args[0]

			i_cfg, err := config.Read(indep_config)
			if err != nil {
				return fmt.Errorf("failed to read config file: %s", err)
			}

			n_cfg_file := i_cfg.Node_Config_Dir + "/config.toml"
			a_cfg_file := i_cfg.Node_Config_Dir + "/app.toml"
			n_cfg, a_cfg, err := config.ReadNodeConfig(n_cfg_file, a_cfg_file)
			if err != nil {
				return fmt.Errorf("failed to read config file: %s", err)
			}

			client, err := client.NewClient(n_cfg.RPC.Address, a_cfg.GRPC.Address)
			if err != nil {

				return fmt.Errorf("failed to conn RPC and GRPC: %s", err)
			}
			defer client.Stop()

			client.GRPC.GetState().String()
			println(client.GRPC.GetState().String())
			println(n_cfg.NODE_NAME)
			println(a_cfg.API.Address)
			tmclient := tmtype.NewServiceClient(client.GRPC.ClientConn)

			LatestBlock, err := tmclient.GetLatestBlock(context.Background(), &tmtype.GetLatestBlockRequest{})
			if err != nil {
				return fmt.Errorf("failed to get LatestBlock: %s", err)
			}
			fmt.Println(LatestBlock.GetBlock().Header.Height)
			LatestValidatorSet, err := tmclient.GetLatestValidatorSet(context.Background(), &tmtype.GetLatestValidatorSetRequest{})
			if err != nil {
				return fmt.Errorf("failed to get LatestValidatorSet: %s", err)
			}
			fmt.Println(append(LatestValidatorSet.GetValidators()))
			return nil
		},
	}
	return cmd
}
