type BPDU struct {
    Version       int               // Protocol version
    SourceSwitchID int               // ID of the sending switch
    RootBridge     int               // ID of the root bridge in the network
    PortStates     map[int]bool      // Port states of the sending switch (e.g., blocked or unblocked)
}
