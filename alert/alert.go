package alert

import (
	"context"
	"fmt"
	"time"

	"github.com/b-harvest/indep_node_alarm_go/client/grpc"
)

func Height_stucked(ctx context.Context, client *grpc.Client, Height_Increasing_Time_Period string) {
	current_block, _ := client.GetLBlock(ctx)
	current_height := current_block.Height
	before_block, _ := client.GetLByBlock(ctx, current_height-1)
	current_block_time_difference := current_block.Time.Sub(before_block.Time)
	for {
		//ParseTime, _ := time.ParseDuration("1s"
		//fmt.Println(ParseTime)
		time.Sleep(current_block_time_difference)
		sync_info, _ := client.GetSyncInfo(ctx)
		if sync_info {
			fmt.Println("This node syncing")
			//alert sync
			// The current block, the block that should eventually be synced, and the average sync rate per block to increase the sleep time.

		} else {
			latest_block, _ := client.GetLBlock(ctx)
			latest_block_height := latest_block.Height
			latest_block_time := latest_block.Time
			latest_block_time_difference := time.Since(latest_block_time)
			block_load_time := latest_block_time_difference.Seconds() - current_block_time_difference.Seconds()
			fmt.Println("block_load: ", block_load_time)
			fmt.Println("latest: ", latest_block_height, "/", latest_block_time_difference.Seconds())
			fmt.Println("curret: ", current_height, "/", current_block_time_difference.Seconds())
			if latest_block_height == current_height || current_block_time_difference.Seconds() < block_load_time {
				alarm_content := "node_name" + ": height stucked!"
				fmt.Println(alarm_content)
				//	alert send_alarm(True, True, alarm_content)

			}
			current_block = latest_block
			current_height = current_block.Height
			before_block, _ := client.GetLByBlock(ctx, current_height-1)
			current_block_time_difference = current_block.Time.Sub(before_block.Time)
		}
	}

}
