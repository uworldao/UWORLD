package dpos

type CandidatesInfo struct {
	Address string
	PeerId  string
}

// initialCandidates the first super node of the block generation cycle.
// The first half is the address of the block, the second half is the id of the block node
var initialCandidates = []CandidatesInfo{
	{
		Address: "UWDNoJGwz4rVACBqfYvQJZxLCSA5o8V5nekV",
		PeerId:  "16Uiu2HAm429SEV84wiyFoCGY1NG5tJbNtqCkMcxsX7haQVhrgaJB",
	},
	{
		Address: "UWDVW6KfmUGjzPiQ36e4LFvPTSo8V2vAh5c3",
		PeerId:  "16Uiu2HAkyd8FzxuGVa3mWiiXthKkAzp2yB6HDeSM7cv5ir8jrmXv",
	},
	{
		Address: "UWDaUrRHiPKRprsD6AYNHiyiaVepZCsbLKYo",
		PeerId:  "16Uiu2HAm58R6j3zpu1nD9rtRoZASA7NSHhybHvoncMECXBDNUA9A",
	},
	{
		Address: "UWDSKmdqJkD8RwiUq2HsiHYwjwvdeKDtTCWv",
		PeerId:  "16Uiu2HAmQASyqisdpSt9HupzGg7JqKXYAU8vL8wfhLcbb3cxR4hs",
	},
	{
		Address: "UWDb62iQKvD4z6QqJW4rYobvLbfmPEBskg5h",
		PeerId:  "16Uiu2HAmGt4s4XmunZi3YwhigwfFuBCExypKjQcwefrJAJrTfPDZ",
	},
	{
		Address: "UWDSyL6uAZvwqhYdxx9kHh4if3w3QbT6f38B",
		PeerId:  "16Uiu2HAm1X2MjfVdj8zq5mXvEyCJ37FkMNFsx6rsWzeWmC8tJjmv",
	},
	{
		Address: "UWDWkzHfwcydwFJEubVpT4XzYHXJoXERuCsZ",
		PeerId:  "16Uiu2HAmVATgQKgF1gb4UdTuABM8i5896fvsCKQzYrqmuR4q7sBg",
	},
	{
		Address: "UWDXiZVnKFknPYwQvsj8Ph7eUXeACnu9ZXPF",
		PeerId:  "16Uiu2HAm8ve2K63xkL3QXixYt7K6BasHFxpvUMeACnGYsXKZ6L9t",
	},
	{
		Address: "UWDYL4adyJmYWryM3r7QPuHduo5q3Tcm31b3",
		PeerId:  "16Uiu2HAmNd8LBBd4tDN6CC9tfXjcMBSDV7cFQgsAadq9ML6fxAam",
	},
	{
		Address: "UWDKoLj4mRTKr4SjyyFG4LY3ExZVSZT9dNZv",
		PeerId:  "16Uiu2HAkziHs8kdy71ASWGqJCa9wfEZ6c4KS5mJk4GVboYeVQA6R",
	},
	{
		Address: "UWDFYhVvEaGW23N57W23iaCbLxcD933jpwsQ",
		PeerId:  "16Uiu2HAkvKciDaQb26ST7VUgVEsEuc2phE7Hff2J49an8t5ZXwPC",
	},
}
