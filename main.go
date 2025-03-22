package main

import "test_component/cmd"

func main() {

	cmd.Execute()

}

// ----------------------------------------------------------------------
//func main() {
//	ctx := context.Background()
//	timeout, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
//	defer cancelFunc()
//	testConn.ContTest(timeout)
//}
