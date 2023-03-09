package cmd

import (
	"context"
	"fmt"

	"github.com/b-harvest/indep_node_alarm_go/client"
	"github.com/b-harvest/indep_node_alarm_go/config"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/privval"
)

func RootCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "indep [config_path]",
		Args:  cobra.ExactArgs(1),
		Short: "Example: $indep ./config.toml",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			indep_config := args[0]

			i_cfg, err := config.Read(indep_config)
			if err != nil {
				return fmt.Errorf("failed to read config file: %s", err)
			}

			n_cfg_file := i_cfg.Node_Home_Dir + "/config/config.toml"
			a_cfg_file := i_cfg.Node_Home_Dir + "/config/app.toml"
			n_cfg, a_cfg, err := config.ReadNodeConfig(n_cfg_file, a_cfg_file)
			if err != nil {
				return fmt.Errorf("failed to read config file: %s", err)
			}

			client, err := client.NewClient(n_cfg.RPC.Address, a_cfg.GRPC.Address)
			if err != nil {
				//alert
				return fmt.Errorf("failed to conn RPC and GRPC: %s", err)
			}
			defer client.Stop()

			//NodeInfo, _ := client.RPC.Client.Status(ctx)
			//NodeInfo.ValidatorInfo.PubKey.Address()
			//page := 1
			//perPage := 1000
			//Validators, _ := client.RPC.Validators(ctx, &NodeInfo.SyncInfo.LatestBlockHeight, &page, &perPage)
			keyFilePath := i_cfg.Node_Home_Dir + "/" + n_cfg.PRIV_VAL_PATH
			stateFilePath := i_cfg.Node_Home_Dir + "/" + n_cfg.PRIV_STATE_PATH
			pv := privval.LoadFilePV(keyFilePath, stateFilePath)
			Validators, _ := client.GRPC.GetLValidatorSet(ctx)
			var validator_runing bool = false

			for _, Validator := range Validators {
				ValidatorAddress := Validator.GetAddress()
				if pv.GetAddress().String() == ValidatorAddress {
					validator_runing = true
					break
				}
			}

			if validator_runing {
				fmt.Println("This node is a validator")
				//missing block func go-rutin
			} else {
				fmt.Println("This node is not a validator")
			}

			//client.RPC.ABCIQuery(ctx, []byte("validator-set"), nil)
			//LatestBlock, _ := client.GRPC.GetLBlock(ctx)
			//if err != nil {
			//	return fmt.Errorf("failed to conn RPC and GRPC: %s", err)
			//}
			//for _, Signature := range LatestBlock.LastCommit.Signatures {
			//	ValidatorAddress := Signature.ValidatorAddress
			//	fmt.Println(ValidatorAddress)
			//}
			return nil
		},
	}
	return cmd
}
