package provider

//------------------- Azure VM Restores Struct ---------------------------



type AzureVMRestorePointsResults struct {
	ID                                       string  `json:"id"`
	BackupDestination                        string  `json:"backupDestination"`
	Type                                     string  `json:"type"`
	VbrID                                    *string `json:"vbrId"`
	PointInTime                              *string `json:"PointInTime,omitempty"`
	PointInTimeLocalTime                     *string `json:"PointInTimeLocalTime,omitempty"`
	BackupSizeBytes                          int     `json:"backupSizeBytes"`
	IsCorrupted                              *bool   `json:"isCorrupted,omitempty"`
	VMName                                   string  `json:"vmName"`
	ResourceHashID                           string  `json:"resourceHashId"`
	RegionID                                 *string `json:"regionId,omitempty"`
	RegionName                               *string `json:"regionName,omitempty"`
	State                                    string  `json:"state"`
	GfsFlags                                 string  `json:"gfsFlags"`
	JobSessionID                             *string `json:"jobSessionId,omitempty"`
	DataRetrievalStatus                      *string `json:DataRetrievalStatus,omitempty`
	RetrievedDataExpirationDate              *string `json:"retrievedDataExpirationDate,omitempty"`
	NotifyBeforeRetrievedDataExpirationHours *int    `json:notifyBeforeRetrievedDataExpirationHours,omitempty`
	ImmutableTill                            *string `json:"immutableTill,omitempty"`
	AccessTier                               *string `json:"accessTier,omitempty"`
	LatestChainSizeBytes                     *int    `json:latestChainSizeBytes,omitempty`
}


type AzureVMRestorePointDataSourceModel struct {
	RestorePointID string `json:"restorePointId"`
}