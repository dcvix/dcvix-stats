//  SPDX-FileCopyrightText: 2025 Diego Cortassa
//  SPDX-License-Identifier: MIT

package globals

const AppName = "Dcvix DCV server stats"

var Metrics = []string{
	"active_streams",
	"stream_sent",
	"stream_recv",
	"dgram_sent",
	"dgram_recv",
	"sent_total_dgrams",
	"recv_total_dgrams",
	"recv_used_dgrams",
	"recv_lost_dgrams",
	"recv_malformed_dgrams",
	"recv_duplicate_dgrams",
	"recv_redundant_dgrams",
	"recv_late_dgrams",
	"recv_dgram_messages_lost",
	"recv_dgram_messages_incomplete",
	"recv_dgram_messages_timegraced",
	"recv_dgram_messages_complete",
	"recv_dgram_messages_inflight",
	"recv_dgram_messages_ready",
	"quic_sent_packets",
	"quic_sent_packets_avg",
	"quic_recv_packets",
	"quic_recv_packets_avg",
	"quic_lost_packets",
	"quic_lost_packets_avg",
	"quic_rtt_nanos",
	"quic_cwnd_size",
	"quic_delivery_rate",
	"intermediates_rtt_nanos",
}

var LogFile string
var LogEntriesQty = 120
var Verbose = false
var RefreshInterval = 30
