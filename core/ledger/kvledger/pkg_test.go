/*
Copyright IBM Corp. 2016 All Rights Reserved.

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

package kvledger

import (
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/spf13/viper"
)

type testEnv struct {
	t    testing.TB
	path string
}

func newTestEnv(t testing.TB) *testEnv {
	path := filepath.Join(
		os.TempDir(),
		"fabric",
		"ledgertests",
		"kvledger",
		strconv.Itoa(rand.Int()))
	return createTestEnv(t, path)
}

func createTestEnv(t testing.TB, path string) *testEnv {
	env := &testEnv{
		t:    t,
		path: path}
	env.cleanup()
	viper.Set("peer.fileSystemPath", env.path)
	return env
}

func (env *testEnv) cleanup() {
	os.RemoveAll(env.path)
}

func minimumOperations(leaves string) int {
	length := len(leaves)
	dp := make([][]int, length)
	for i := range dp {
		dp[i] = make([]int, 3)
	}
	// dp[i][j]: 前i个树叶处于状态j的最小操作数
	if leaves[0] == 'r' {
		dp[0][0] = 0
	}else {
		dp[0][0] = 1
	}
	dp[0][1] = math.MaxInt32
	dp[0][2] = math.MaxInt32
	dp[1][2] = math.MaxInt32

	isYellow := -1
	if leaves[1] == 'y' {
		isYellow = 0
	}else {
		isYellow = 1
	}
	dp[1][0] = dp[0][0] + (1-isYellow)
	dp[1][1] = min(dp[0][0], dp[0][1]) + isYellow
	dp[1][2] = min(dp[0][1], dp[0][2]) + (1 - isYellow)

	for i := 2; i < length; i++ {
		if leaves[1] == 'y' {
			isYellow = 0
		}else {
			isYellow = 1
		}
		dp[i][0] = dp[i-1][0] + (1-isYellow)
		dp[i][1] = min(dp[i-1][0], dp[i-1][1]) + isYellow
		dp[i][2] = min(dp[i-1][1], dp[i-1][2]) + (1-isYellow)
	}
	return dp[length-1][2]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


