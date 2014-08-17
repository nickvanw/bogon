package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

const command = "ohai"

func init() {
	AddPlugin("Sysinfo", "(?i)^\\.sysinfo$", MessageHandler(SysInfo), false, true)
}

func SysInfo(msg *Message) {
	var OhaiOutput OhaiData
	binary, lookErr := exec.LookPath(command)
	if lookErr != nil {
		msg.Return("The binary ohai is not installed, sorry")
		return
	}
	cmd := exec.Command(binary)
	data, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing `ohai`")
		return
	}
	json.Unmarshal(data, &OhaiOutput)
	output_str := make([]string, 0)
	os_string := fmt.Sprintf("%s %s %s (%s) | Uptime: %s",
		OhaiOutput.Os, OhaiOutput.Platform, OhaiOutput.PlatformVersion, OhaiOutput.OsVersion, OhaiOutput.Uptime)
	output_str = append(output_str, os_string)
	cpu_data, cpu_data_ok := OhaiOutput.Cpu["0"]
	var cpu_string string
	if cpu_data_ok {
		cpu_cores, cpu_cores_ok := OhaiOutput.Cpu["total"]
		if cpu_cores_ok {
			cpu_string = fmt.Sprintf("%v x ", cpu_cores.(float64))
		}
		cpu_name := cpu_data.(map[string]interface{})["model_name"].(string)
		cpu_string = cpu_string + fmt.Sprintf("%s", cpu_name)
	}
	output_str = append(output_str, cpu_string)
	r := strings.NewReplacer("kB", "")
	mt, mt_er := strconv.ParseInt(r.Replace(OhaiOutput.Memory.Total), 0, 64)
	mf, mf_er := strconv.ParseInt(r.Replace(OhaiOutput.Memory.Free), 0, 64)
	mst, mst_er := strconv.ParseInt(r.Replace(OhaiOutput.Memory.Swap.Total), 0, 64)
	msf, msf_er := strconv.ParseInt(r.Replace(OhaiOutput.Memory.Swap.Free), 0, 64)
	var mem_string string
	if mt_er == nil && mf_er == nil && mst_er == nil && msf_er == nil {
		mem_percent := (1.0 - (float64(mf) / float64(mt))) * 100
		swp_percent := (1.0 - (float64(msf) / float64(mst))) * 100
		mem_used := uint64(mt - mf)
		swp_used := uint64(mst - msf)
		mem_string = fmt.Sprintf("Mem: %s Total, %s Used (%.2f%%) | Swap: %s Total, %s Used (%.2f%%)",
			humanize.Bytes(uint64(mt)*1024), humanize.Bytes(mem_used*1024), mem_percent,
			humanize.Bytes(uint64(mst)*2014), humanize.Bytes(uint64(swp_used)*1024), swp_percent)
	}
	output_str = append(output_str, mem_string)
	var nic_info string
	def_in, def_in_ok := OhaiOutput.Counters.Network.Interfaces[OhaiOutput.Network.DefaultInterface]
	if def_in_ok {
		tx, tx_e := strconv.ParseInt(def_in.TX.Bytes, 0, 64)
		rx, rx_e := strconv.ParseInt(def_in.RX.Bytes, 0, 64)
		if tx_e == nil && rx_e == nil {
			transmit := humanize.Bytes(uint64(tx))
			recieve := humanize.Bytes(uint64(rx))
			nic_info = fmt.Sprintf("Net: %s: %s tx / %s rx", OhaiOutput.Network.DefaultInterface,
				transmit, recieve)
		}
	}
	output_str = append(output_str, nic_info)
	msg.Return(strings.Join(output_str, " | "))
}

type OhaiData struct {
	Counters struct {
		Network struct {
			Interfaces map[string]NetworkInterface `json:"interfaces"`
		} `json:"network"`
	} `json:"counters"`
	Cpu         map[string]interface{} `json:"cpu"`
	CurrentUser string                 `json:"current_user"`
	Hostname    string                 `json:"hostname"`
	Ipaddress   string                 `json:"ipaddress"`
	Network     struct {
		DefaultInterface string `json:"default_interface"`
	} `json:"network"`
	Memory struct {
		Active          string `json:"active"`
		AnonPages       string `json:"anon_pages"`
		Buffers         string `json:"buffers"`
		Cached          string `json:"cached"`
		Dirty           string `json:"dirty"`
		Free            string `json:"free"`
		Inactive        string `json:"inactive"`
		Slab            string `json:"slab"`
		SlabReclaimable string `json:"slab_reclaimable"`
		SlabUnreclaim   string `json:"slab_unreclaim"`
		Swap            struct {
			Free  string `json:"free"`
			Total string `json:"total"`
		} `json:"swap"`
		Total     string `json:"total"`
		Writeback string `json:"writeback"`
	} `json:"memory"`
	OhaiTime        float64 `json:"ohai_time"`
	Os              string  `json:"os"`
	OsVersion       string  `json:"os_version"`
	Platform        string  `json:"platform"`
	PlatformBuild   string  `json:"platform_build"`
	PlatformFamily  string  `json:"platform_family"`
	PlatformVersion string  `json:"platform_version"`
	Uptime          string  `json:"uptime"`
	UptimeSeconds   float64 `json:"uptime_seconds"`
}

type NetworkInterface struct {
	RX struct {
		Bytes      string  `json:"bytes"`
		Carrier    float64 `json:"carrier"`
		Collisions string  `json:"collisions"`
		Compressed float64 `json:"compressed"`
		Drop       float64 `json:"drop"`
		Errors     string  `json:"errors"`
		Overrun    float64 `json:"overrun"`
		Packets    string  `json:"packets"`
	} `json:"rx"`
	TX struct {
		Bytes      string  `json:"bytes"`
		Carrier    float64 `json:"carrier"`
		Collisions string  `json:"collisions"`
		Compressed float64 `json:"compressed"`
		Drop       float64 `json:"drop"`
		Errors     string  `json:"errors"`
		Overrun    float64 `json:"overrun"`
		Packets    string  `json:"packets"`
	} `json:"tx"`
}
