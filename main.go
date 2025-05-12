package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/autobrr/go-qbittorrent"
	"gopkg.in/yaml.v3"
)

var STATIC_TRACKERS = []string{
	"udp://tracker.coppersurfer.tk:6969/announce",
	"http://tracker.internetwarriors.net:1337/announce",
	"udp://tracker.internetwarriors.net:1337/announce",
	"udp://tracker.opentrackr.org:1337/announce",
	"udp://9.rarbg.to:2710/announce",
	"udp://exodus.desync.com:6969/announce",
	"udp://explodie.org:6969/announce",
	"http://explodie.org:6969/announce",
	"udp://public.popcorn-tracker.org:6969/announce",
	"udp://tracker.vanitycore.co:6969/announce",
	"http://tracker.vanitycore.co:6969/announce",
	"udp://tracker1.itzmx.com:8080/announce",
	"http://tracker1.itzmx.com:8080/announce",
	"udp://ipv4.tracker.harry.lu:80/announce",
	"udp://tracker.torrent.eu.org:451/announce",
	"udp://tracker.tiny-vps.com:6969/announce",
	"udp://tracker.port443.xyz:6969/announce",
	"udp://open.stealth.si:80/announce",
	"udp://open.demonii.si:1337/announce",
	"udp://denis.stalker.upeer.me:6969/announce",
	"udp://bt.xxx-tracker.com:2710/announce",
	"http://tracker.port443.xyz:6969/announce",
	"udp://tracker2.itzmx.com:6961/announce",
	"udp://retracker.lanta-net.ru:2710/announce",
	"http://tracker2.itzmx.com:6961/announce",
	"http://tracker4.itzmx.com:2710/announce",
	"http://tracker3.itzmx.com:6961/announce",
	"http://tracker.city9x.com:2710/announce",
	"http://torrent.nwps.ws:80/announce",
	"http://retracker.telecom.by:80/announce",
	"http://open.acgnxtracker.com:80/announce",
	"wss://ltrackr.iamhansen.xyz:443/announce",
	"udp://zephir.monocul.us:6969/announce",
	"udp://tracker.toss.li:6969/announce",
	"http://opentracker.xyz:80/announce",
	"http://open.trackerlist.xyz:80/announce",
	"udp://tracker.swateam.org.uk:2710/announce",
	"udp://tracker.kamigami.org:2710/announce",
	"udp://tracker.iamhansen.xyz:2000/announce",
	"udp://tracker.ds.is:6969/announce",
	"udp://pubt.in:2710/announce",
	"https://tracker.fastdownload.xyz:443/announce",
	"https://opentracker.xyz:443/announce",
	"http://tracker.torrentyorg.pl:80/announce",
	"http://t.nyaatracker.com:80/announce",
	"http://open.acgtracker.com:1096/announce",
	"wss://tracker.openwebtorrent.com:443/announce",
	"wss://tracker.fastcast.nz:443/announce",
	"wss://tracker.btorrent.xyz:443/announce",
	"udp://tracker.justseed.it:1337/announce",
	"udp://thetracker.org:80/announce",
	"udp://packages.crunchbangplusplus.org:6969/announce",
	"https://1337.abcvg.info:443/announce",
	"http://tracker.tfile.me:80/announce.php",
	"http://tracker.tfile.me:80/announce",
	"http://tracker.tfile.co:80/announce",
	"http://retracker.mgts.by:80/announce",
	"http://peersteers.org:80/announce",
	"http://fxtt.ru:80/announce",
}

var STATIC_TRACKER_URLS = []string{
	"https://newtrackon.com/api/stable?include_ipv4_only_trackers=true&include_ipv6_only_trackers=false",
	"https://trackerslist.com/best.txt",
	"https://trackerslist.com/http.txt",
	"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_best.txt",
	"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_https.txt",
	"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all_i2p.txt",
}

// Constants
const (
	NEW_TRACKON_API_URL = "https://newtrackon.com/api/add"
)

// Config holds all configuration for the application.
type Config struct {
	Host                  string   `json:"host" yaml:"host" toml:"host"`
	Username              string   `json:"username" yaml:"username" toml:"username"`
	Password              string   `json:"password" yaml:"password" toml:"password"`
	TLSSkipVerify         bool     `json:"tlsSkipVerify" yaml:"tlsSkipVerify" toml:"tlsSkipVerify"`
	ContributeTrackers    bool     `json:"contributeTrackers" yaml:"contributeTrackers" toml:"contributeTrackers"`
	ReannounceTorrents    bool     `json:"reannounceTorrents" yaml:"reannounceTorrents" toml:"reannounceTorrents"`
	StaticTrackers        []string `json:"staticTrackers" yaml:"staticTrackers" toml:"staticTrackers"`
	TrackerListURLs       []string `json:"trackerListURLs" yaml:"trackerListURLs" toml:"trackerListURLs"`
	ReannounceMaxAttempts int      `json:"reannounceMaxAttempts" yaml:"reannounceMaxAttempts" toml:"reannounceMaxAttempts"`
	ReannounceInterval    int      `json:"reannounceInterval" yaml:"reannounceInterval" toml:"reannounceInterval"`
}

// loadConfigFromFile loads configuration from a file in YAML, TOML, or JSON format.
func loadConfigFromFile(filepath string) (*Config, error) {
	fmt.Printf("INFO: Loading configuration from file: %s", filepath)

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	cfg := &Config{}

	// Try to determine the format based on file extension
	ext := strings.ToLower(filepath[strings.LastIndexByte(filepath, '.')+1:])
	switch ext {
	case "yaml", "yml":
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("error parsing YAML config: %w", err)
		}
		fmt.Println("INFO: Parsed YAML configuration file")
	case "toml":
		if err := toml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("error parsing TOML config: %w", err)
		}
		fmt.Println("INFO: Parsed TOML configuration file")
	case "json":
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("error parsing JSON config: %w", err)
		}
		fmt.Println("INFO: Parsed JSON configuration file")
	default:
		// Try each format in order: YAML, TOML, JSON
		if err := yaml.Unmarshal(data, cfg); err == nil {
			fmt.Println("INFO: Parsed YAML configuration file")
			return cfg, nil
		}

		if err := toml.Unmarshal(data, cfg); err == nil {
			fmt.Println("INFO: Parsed TOML configuration file")
			return cfg, nil
		}

		if err := json.Unmarshal(data, cfg); err == nil {
			fmt.Println("INFO: Parsed JSON configuration file")
			return cfg, nil
		}

		return nil, fmt.Errorf("unrecognized config file format for %s", filepath)
	}

	return cfg, nil
}

// tryDefaultConfigFiles tries to load configuration from default config files in the order:
// 1. config.yaml or config.yml
// 2. config.toml
// 3. config.json
func tryDefaultConfigFiles() (*Config, bool) {
	defaultPaths := []string{
		"config.yaml",
		"config.yml",
		"config.toml",
		"config.json",
	}

	for _, path := range defaultPaths {
		if _, err := os.Stat(path); err == nil {
			cfg, err := loadConfigFromFile(path)
			if err == nil {
				return cfg, true
			}
			fmt.Printf("WARN: Found default config file %s but could not parse it: %v", path, err)
		}
	}

	return nil, false
}

// loadConfig loads configuration from CLI flags, environment variables, and config files.
// Order of precedence (highest to lowest):
// 1. CLI flags
// 2. Environment variables
// 3. Config file specified by --config flag or QBIT_CONFIG_FILE environment variable
// 4. Default config files (config.yaml/yml, config.toml, config.json)
func loadConfig() (*Config, error) {
	// Create a new FlagSet just for the config file path to avoid double-parsing issues
	configFlags := flag.NewFlagSet("configFlags", flag.ContinueOnError)
	configFlags.SetOutput(io.Discard) // Don't print usage errors for this stage
	configPath := configFlags.String("config", "", "Path to config file (YAML, TOML, or JSON)")

	// Try to parse only the config flag first
	_ = configFlags.Parse(os.Args[1:])

	// Also check environment variable for config path
	cfgFile := *configPath
	if cfgFile == "" {
		cfgFile = os.Getenv("QBIT_CONFIG_FILE")
	}

	// Try to load from config file if specified
	var cfg *Config
	var loadedFromFile bool

	if cfgFile != "" {
		var err error
		cfg, err = loadConfigFromFile(cfgFile)
		if err != nil {
			return nil, fmt.Errorf("could not load config file %s: %w", cfgFile, err)
		}
		loadedFromFile = true
	} else {
		// Try default config files
		cfg, loadedFromFile = tryDefaultConfigFiles()
		if !loadedFromFile {
			cfg = &Config{} // Start with empty config if no file was found
		}
	}

	// Define all CLI flags, including the config flag (so it shows up in usage)
	flag.StringVar(&cfg.Host, "host", defaultValueOrEnv(cfg.Host, os.Getenv("QBIT_HOST"), ""), "qBittorrent host (e.g., http://localhost:8080). Env: QBIT_HOST")
	flag.StringVar(&cfg.Username, "username", defaultValueOrEnv(cfg.Username, os.Getenv("QBIT_USERNAME"), ""), "qBittorrent username. Env: QBIT_USERNAME")
	flag.StringVar(&cfg.Password, "password", defaultValueOrEnv(cfg.Password, os.Getenv("QBIT_PASSWORD"), ""), "qBittorrent password. Env: QBIT_PASSWORD")
	flag.BoolVar(&cfg.TLSSkipVerify, "tlsSkipVerify", defaultBoolValueOrEnv(cfg.TLSSkipVerify, os.Getenv("QBIT_TLS_SKIP_VERIFY"), false), "Skip TLS certificate verification for qBittorrent client. Env: QBIT_TLS_SKIP_VERIFY")
	flag.BoolVar(&cfg.ContributeTrackers, "contributeTrackers", defaultBoolValueOrEnv(cfg.ContributeTrackers, os.Getenv("QBIT_CONTRIBUTE_TRACKERS"), false), "Contribute found trackers to newtrackon.com. Env: QBIT_CONTRIBUTE_TRACKERS")
	flag.BoolVar(&cfg.ReannounceTorrents, "reannounceTorrents", defaultBoolValueOrEnv(cfg.ReannounceTorrents, os.Getenv("QBIT_REANNOUNCE_TORRENTS"), false), "Reannounce all torrents after updating trackers. Env: QBIT_REANNOUNCE_TORRENTS")
	configFileFlag := flag.String("config", cfgFile, "Path to config file (YAML, TOML, or JSON). Env: QBIT_CONFIG_FILE")

	// Parse all flags
	flag.Parse()

	// Update config file path if changed during main flag parse
	if *configFileFlag != cfgFile && *configFileFlag != "" {
		// If config flag was provided during main flag parsing and different from what we initially found,
		// reload config from file
		newCfg, err := loadConfigFromFile(*configFileFlag)
		if err != nil {
			return nil, fmt.Errorf("could not load config file %s: %w", *configFileFlag, err)
		}

		// Only use values from file that weren't explicitly set via command line
		flag.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "host":
				newCfg.Host = cfg.Host
			case "username":
				newCfg.Username = cfg.Username
			case "password":
				newCfg.Password = cfg.Password
			case "tlsSkipVerify":
				newCfg.TLSSkipVerify = cfg.TLSSkipVerify
			case "contributeTrackers":
				newCfg.ContributeTrackers = cfg.ContributeTrackers
			case "reannounceTorrents":
				newCfg.ReannounceTorrents = cfg.ReannounceTorrents
			}
		})

		cfg = newCfg
	}

	// Set defaults for fields that weren't loaded or specified
	if len(cfg.StaticTrackers) == 0 {
		cfg.StaticTrackers = STATIC_TRACKERS
	}

	if len(cfg.TrackerListURLs) == 0 {
		cfg.TrackerListURLs = STATIC_TRACKER_URLS
	}

	if cfg.ReannounceMaxAttempts == 0 {
		cfg.ReannounceMaxAttempts = 3
	}

	if cfg.ReannounceInterval == 0 {
		cfg.ReannounceInterval = 5
	}

	if cfg.Host == "" {
		return nil, fmt.Errorf("qBittorrent host must be set via --host flag, QBIT_HOST environment variable, or in config file")
	}

	return cfg, nil
}

// Helper functions for default values with ENV fallback
func defaultValueOrEnv(configVal, envVal, defaultVal string) string {
	if configVal != "" {
		return configVal
	}
	if envVal != "" {
		return envVal
	}
	return defaultVal
}

func defaultBoolValueOrEnv(configVal bool, envVal string, defaultVal bool) bool {
	if envVal != "" {
		parsedVal, err := strconv.ParseBool(envVal)
		if err == nil {
			return parsedVal
		}
		fmt.Printf("WARN: Invalid boolean value for environment variable: '%s'. Using default: %t", envVal, defaultVal)
	}
	return configVal
}

func getTrackersFromUrl(cfg *Config, url string) (*[]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(bodyBytes)
	trackers := parseTrackers(cfg, bodyString)

	return &trackers, nil
}

func getTrackerLists(cfg *Config) (*[]string, error) {
	trackers := make([]string, 0)

	for _, trackerUrl := range cfg.TrackerListURLs {
		trackerList, err := getTrackersFromUrl(cfg, trackerUrl)
		if err != nil {
			return nil, err
		}

		trackers = append(trackers, *trackerList...)
	}

	return &trackers, nil
}

func parseTrackers(cfg *Config, s string) []string {
	stringParts := strings.Split(s, "\n")

	trackers := make([]string, 0)

	// Add static trackers from config
	trackers = append(trackers, cfg.StaticTrackers...)

	for _, tracker := range stringParts {
		trimmedTracker := strings.TrimSpace(tracker)
		if trimmedTracker != "" {
			trackers = append(trackers, trimmedTracker)
		}
	}

	// Dedupe the list of trackers
	trackers = removeDuplicate(trackers)

	return trackers
}

func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func submitTrackersToNewTrackon(trackers []string) error {
	apiURL := NEW_TRACKON_API_URL
	data := url.Values{}
	// Join with space as per newtrackon.com API examples
	data.Set("new_trackers", strings.Join(trackers, " "))

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// newtrackon.com API docs specify 204 No Content on success.
	// We also check for 200 OK as a general success status.
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("WARN: Failed to submit trackers to newtrackon.com. Status: %s, Body: %s", resp.Status, string(bodyBytes))
		// Not returning an error here as it's a non-critical part of the script's main function
	} else {
		fmt.Println("INFO: Successfully submitted trackers to newtrackon.com")
	}

	return nil
}

// runApp contains the core logic of the application.
func runApp(cfg *Config) error {
	fmt.Printf("INFO: Connecting to qBittorrent at %s (TLS Skip Verify: %t)", cfg.Host, cfg.TLSSkipVerify)

	qclient := qbittorrent.NewClient(qbittorrent.Config{
		Host:          cfg.Host,
		Username:      cfg.Username,
		Password:      cfg.Password,
		TLSSkipVerify: cfg.TLSSkipVerify,
	})

	ctx := context.Background()

	if err := qclient.LoginCtx(ctx); err != nil {
		return fmt.Errorf("could not log into qBittorrent client: %w", err)
	}
	fmt.Println("INFO: Successfully logged into qBittorrent client")

	torrents, err := qclient.GetTorrents(qbittorrent.TorrentFilterOptions{})
	if err != nil {
		return fmt.Errorf("could not get torrents from client: %w", err)
	}
	fmt.Printf("INFO: Found %d torrents", len(torrents))

	fetchedTrackers, err := getTrackerLists(cfg)
	if err != nil {
		return fmt.Errorf("could not get trackers list: %w", err)
	}
	fmt.Printf("INFO: Fetched %d unique trackers", len(*fetchedTrackers))

	if cfg.ContributeTrackers && len(*fetchedTrackers) > 0 {
		fmt.Println("INFO: Contributing trackers to newtrackon.com.")
		if err := submitTrackersToNewTrackon(*fetchedTrackers); err != nil {
			// Log the error but don't stop the main script execution as it's non-critical
			fmt.Printf("WARN: Could not submit trackers to newtrackon.com: %v", err)
		}
	}

	fmt.Println("INFO: Updating qBittorrent client's default trackers...")
	trackerStrings := strings.Join(*fetchedTrackers, "\n\n")

	if err = qclient.SetPreferences(map[string]interface{}{
		"add_trackers_enabled": true,
		"add_trackers":         trackerStrings,
	}); err != nil {
		return fmt.Errorf("could not update qBittorrent client's default trackers: %w", err)
	}
	fmt.Println("INFO: Successfully updated qBittorrent client's default trackers.")

	fmt.Println("INFO: Updating trackers for existing torrents...")
	updatedCount := 0
	skippedCount := 0
	for _, torrent := range torrents {
		if err = qclient.AddTrackers(torrent.Hash, trackerStrings); err != nil {
			fmt.Printf("WARN: Could not update trackers for torrent %s (%s): %v", torrent.Name, torrent.Hash, err)
			skippedCount++
		} else {
			updatedCount++
		}
	}
	fmt.Printf("INFO: Finished updating trackers for existing torrents. Updated: %d, Skipped/Failed: %d", updatedCount, skippedCount)

	if cfg.ReannounceTorrents {
		fmt.Println("INFO: Reannouncing all torrents.")
		if err = qclient.ReannounceTorrentWithRetry(context.Background(), "all", &qbittorrent.ReannounceOptions{
			DeleteOnFailure: false,
			MaxAttempts:     cfg.ReannounceMaxAttempts,
			Interval:        cfg.ReannounceInterval,
		}); err != nil {
			fmt.Printf("WARN: Failed to reannounce all torrents: %v", err) // Non-fatal, log and continue
		} else {
			fmt.Println("INFO: Successfully initiated reannounce for all torrents.")
		}
	} else {
		fmt.Println("INFO: Skipping reannounce of torrents because it is disabled.")
	}

	fmt.Println("INFO: Script finished.")
	return nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("ERROR: Configuration error: %v", err)
	}

	if err := runApp(cfg); err != nil {
		log.Fatalf("ERROR: Application error: %v", err)
	}
}
