// Code based on the AtlasAPI V2 OpenAPI file

package admin

// GroupSettings Collection of settings that configures the project.
type GroupSettings struct {
	// Flag that indicates whether the AI Cluster Assistant is enabled for the specified project.
	IsClusterAiAssistantEnabled *bool `json:"isClusterAiAssistantEnabled,omitempty"`
	// Flag that indicates whether to collect database-specific metrics for the specified project.
	IsCollectDatabaseSpecificsStatisticsEnabled *bool `json:"isCollectDatabaseSpecificsStatisticsEnabled,omitempty"`
	// Flag that indicates whether to enable the Data Explorer for the specified project.
	IsDataExplorerEnabled *bool `json:"isDataExplorerEnabled,omitempty"`
	// Flag that indicates whether to enable the use of generative AI features which make requests to 3rd party services in Data Explorer for the specified project.
	IsDataExplorerGenAIFeaturesEnabled *bool `json:"isDataExplorerGenAIFeaturesEnabled,omitempty"`
	// Flag that indicates whether to enable the passing of sample field values with the use of generative AI features in the Data Explorer for the specified project.
	IsDataExplorerGenAISampleDocumentPassingEnabled *bool `json:"isDataExplorerGenAISampleDocumentPassingEnabled,omitempty"`
	// Flag that indicates whether to enable extended storage sizes for the specified project.
	IsExtendedStorageSizesEnabled *bool `json:"isExtendedStorageSizesEnabled,omitempty"`
	// Flag that indicates whether to enable Native Reranking with Voyage AI models in the Aggregation Pipeline for the specified project.
	IsNativeRerankingEnabled *bool `json:"isNativeRerankingEnabled,omitempty"`
	// Flag that indicates whether to enable the Performance Advisor and Profiler for the specified project.
	IsPerformanceAdvisorEnabled *bool `json:"isPerformanceAdvisorEnabled,omitempty"`
	// Flag that indicates whether to enable the Real Time Performance Panel for the specified project.
	IsRealtimePerformancePanelEnabled *bool `json:"isRealtimePerformancePanelEnabled,omitempty"`
	// Flag that indicates whether to enable the Schema Advisor for the specified project.
	IsSchemaAdvisorEnabled *bool `json:"isSchemaAdvisorEnabled,omitempty"`
}

// NewGroupSettings instantiates a new GroupSettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupSettings() *GroupSettings {
	this := GroupSettings{}
	var isDataExplorerGenAISampleDocumentPassingEnabled bool = false
	this.IsDataExplorerGenAISampleDocumentPassingEnabled = &isDataExplorerGenAISampleDocumentPassingEnabled
	return &this
}

// NewGroupSettingsWithDefaults instantiates a new GroupSettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupSettingsWithDefaults() *GroupSettings {
	this := GroupSettings{}
	var isDataExplorerGenAISampleDocumentPassingEnabled bool = false
	this.IsDataExplorerGenAISampleDocumentPassingEnabled = &isDataExplorerGenAISampleDocumentPassingEnabled
	return &this
}

// GetIsClusterAiAssistantEnabled returns the IsClusterAiAssistantEnabled field value if set, zero value otherwise
func (o *GroupSettings) GetIsClusterAiAssistantEnabled() bool {
	if o == nil || IsNil(o.IsClusterAiAssistantEnabled) {
		var ret bool
		return ret
	}
	return *o.IsClusterAiAssistantEnabled
}

// GetIsClusterAiAssistantEnabledOk returns a tuple with the IsClusterAiAssistantEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupSettings) GetIsClusterAiAssistantEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.IsClusterAiAssistantEnabled) {
		return nil, false
	}

	return o.IsClusterAiAssistantEnabled, true
}

// HasIsClusterAiAssistantEnabled returns a boolean if a field has been set.
func (o *GroupSettings) HasIsClusterAiAssistantEnabled() bool {
	if o != nil && !IsNil(o.IsClusterAiAssistantEnabled) {
		return true
	}

	return false
}

// SetIsClusterAiAssistantEnabled gets a reference to the given bool and assigns it to the IsClusterAiAssistantEnabled field.
func (o *GroupSettings) SetIsClusterAiAssistantEnabled(v bool) {
	o.IsClusterAiAssistantEnabled = &v
}

// GetIsCollectDatabaseSpecificsStatisticsEnabled returns the IsCollectDatabaseSpecificsStatisticsEnabled field value if set, zero value otherwise
func (o *GroupSettings) GetIsCollectDatabaseSpecificsStatisticsEnabled() bool {
	if o == nil || IsNil(o.IsCollectDatabaseSpecificsStatisticsEnabled) {
		var ret bool
		return ret
	}
	return *o.IsCollectDatabaseSpecificsStatisticsEnabled
}

// GetIsCollectDatabaseSpecificsStatisticsEnabledOk returns a tuple with the IsCollectDatabaseSpecificsStatisticsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupSettings) GetIsCollectDatabaseSpecificsStatisticsEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.IsCollectDatabaseSpecificsStatisticsEnabled) {
		return nil, false
	}

	return o.IsCollectDatabaseSpecificsStatisticsEnabled, true
}

// HasIsCollectDatabaseSpecificsStatisticsEnabled returns a boolean if a field has been set.
func (o *GroupSettings) HasIsCollectDatabaseSpecificsStatisticsEnabled() bool {
	if o != nil && !IsNil(o.IsCollectDatabaseSpecificsStatisticsEnabled) {
		return true
	}

	return false
}

// SetIsCollectDatabaseSpecificsStatisticsEnabled gets a reference to the given bool and assigns it to the IsCollectDatabaseSpecificsStatisticsEnabled field.
func (o *GroupSettings) SetIsCollectDatabaseSpecificsStatisticsEnabled(v bool) {
	o.IsCollectDatabaseSpecificsStatisticsEnabled = &v
}

// GetIsDataExplorerEnabled returns the IsDataExplorerEnabled field value if set, zero value otherwise
func (o *GroupSettings) GetIsDataExplorerEnabled() bool {
	if o == nil || IsNil(o.IsDataExplorerEnabled) {
		var ret bool
		return ret
	}
	return *o.IsDataExplorerEnabled
}

// GetIsDataExplorerEnabledOk returns a tuple with the IsDataExplorerEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupSettings) GetIsDataExplorerEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.IsDataExplorerEnabled) {
		return nil, false
	}

	return o.IsDataExplorerEnabled, true
}

// HasIsDataExplorerEnabled returns a boolean if a field has been set.
func (o *GroupSettings) HasIsDataExplorerEnabled() bool {
	if o != nil && !IsNil(o.IsDataExplorerEnabled) {
		return true
	}

	return false
}

// SetIsDataExplorerEnabled gets a reference to the given bool and assigns it to the IsDataExplorerEnabled field.
func (o *GroupSettings) SetIsDataExplorerEnabled(v bool) {
	o.IsDataExplorerEnabled = &v
}

// GetIsDataExplorerGenAIFeaturesEnabled returns the IsDataExplorerGenAIFeaturesEnabled field value if set, zero value otherwise
func (o *GroupSettings) GetIsDataExplorerGenAIFeaturesEnabled() bool {
	if o == nil || IsNil(o.IsDataExplorerGenAIFeaturesEnabled) {
		var ret bool
		return ret
	}
	return *o.IsDataExplorerGenAIFeaturesEnabled
}

// GetIsDataExplorerGenAIFeaturesEnabledOk returns a tuple with the IsDataExplorerGenAIFeaturesEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupSettings) GetIsDataExplorerGenAIFeaturesEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.IsDataExplorerGenAIFeaturesEnabled) {
		return nil, false
	}

	return o.IsDataExplorerGenAIFeaturesEnabled, true
}

// HasIsDataExplorerGenAIFeaturesEnabled returns a boolean if a field has been set.
func (o *GroupSettings) HasIsDataExplorerGenAIFeaturesEnabled() bool {
	if o != nil && !IsNil(o.IsDataExplorerGenAIFeaturesEnabled) {
		return true
	}

	return false
}

// SetIsDataExplorerGenAIFeaturesEnabled gets a reference to the given bool and assigns it to the IsDataExplorerGenAIFeaturesEnabled field.
func (o *GroupSettings) SetIsDataExplorerGenAIFeaturesEnabled(v bool) {
	o.IsDataExplorerGenAIFeaturesEnabled = &v
}

// GetIsDataExplorerGenAISampleDocumentPassingEnabled returns the IsDataExplorerGenAISampleDocumentPassingEnabled field value if set, zero value otherwise
func (o *GroupSettings) GetIsDataExplorerGenAISampleDocumentPassingEnabled() bool {
	if o == nil || IsNil(o.IsDataExplorerGenAISampleDocumentPassingEnabled) {
		var ret bool
		return ret
	}
	return *o.IsDataExplorerGenAISampleDocumentPassingEnabled
}

// GetIsDataExplorerGenAISampleDocumentPassingEnabledOk returns a tuple with the IsDataExplorerGenAISampleDocumentPassingEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupSettings) GetIsDataExplorerGenAISampleDocumentPassingEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.IsDataExplorerGenAISampleDocumentPassingEnabled) {
		return nil, false
	}

	return o.IsDataExplorerGenAISampleDocumentPassingEnabled, true
}

// HasIsDataExplorerGenAISampleDocumentPassingEnabled returns a boolean if a field has been set.
func (o *GroupSettings) HasIsDataExplorerGenAISampleDocumentPassingEnabled() bool {
	if o != nil && !IsNil(o.IsDataExplorerGenAISampleDocumentPassingEnabled) {
		return true
	}

	return false
}

// SetIsDataExplorerGenAISampleDocumentPassingEnabled gets a reference to the given bool and assigns it to the IsDataExplorerGenAISampleDocumentPassingEnabled field.
func (o *GroupSettings) SetIsDataExplorerGenAISampleDocumentPassingEnabled(v bool) {
	o.IsDataExplorerGenAISampleDocumentPassingEnabled = &v
}

// GetIsExtendedStorageSizesEnabled returns the IsExtendedStorageSizesEnabled field value if set, zero value otherwise
func (o *GroupSettings) GetIsExtendedStorageSizesEnabled() bool {
	if o == nil || IsNil(o.IsExtendedStorageSizesEnabled) {
		var ret bool
		return ret
	}
	return *o.IsExtendedStorageSizesEnabled
}

// GetIsExtendedStorageSizesEnabledOk returns a tuple with the IsExtendedStorageSizesEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupSettings) GetIsExtendedStorageSizesEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.IsExtendedStorageSizesEnabled) {
		return nil, false
	}

	return o.IsExtendedStorageSizesEnabled, true
}

// HasIsExtendedStorageSizesEnabled returns a boolean if a field has been set.
func (o *GroupSettings) HasIsExtendedStorageSizesEnabled() bool {
	if o != nil && !IsNil(o.IsExtendedStorageSizesEnabled) {
		return true
	}

	return false
}

// SetIsExtendedStorageSizesEnabled gets a reference to the given bool and assigns it to the IsExtendedStorageSizesEnabled field.
func (o *GroupSettings) SetIsExtendedStorageSizesEnabled(v bool) {
	o.IsExtendedStorageSizesEnabled = &v
}

// GetIsNativeRerankingEnabled returns the IsNativeRerankingEnabled field value if set, zero value otherwise
func (o *GroupSettings) GetIsNativeRerankingEnabled() bool {
	if o == nil || IsNil(o.IsNativeRerankingEnabled) {
		var ret bool
		return ret
	}
	return *o.IsNativeRerankingEnabled
}

// GetIsNativeRerankingEnabledOk returns a tuple with the IsNativeRerankingEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupSettings) GetIsNativeRerankingEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.IsNativeRerankingEnabled) {
		return nil, false
	}

	return o.IsNativeRerankingEnabled, true
}

// HasIsNativeRerankingEnabled returns a boolean if a field has been set.
func (o *GroupSettings) HasIsNativeRerankingEnabled() bool {
	if o != nil && !IsNil(o.IsNativeRerankingEnabled) {
		return true
	}

	return false
}

// SetIsNativeRerankingEnabled gets a reference to the given bool and assigns it to the IsNativeRerankingEnabled field.
func (o *GroupSettings) SetIsNativeRerankingEnabled(v bool) {
	o.IsNativeRerankingEnabled = &v
}

// GetIsPerformanceAdvisorEnabled returns the IsPerformanceAdvisorEnabled field value if set, zero value otherwise
func (o *GroupSettings) GetIsPerformanceAdvisorEnabled() bool {
	if o == nil || IsNil(o.IsPerformanceAdvisorEnabled) {
		var ret bool
		return ret
	}
	return *o.IsPerformanceAdvisorEnabled
}

// GetIsPerformanceAdvisorEnabledOk returns a tuple with the IsPerformanceAdvisorEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupSettings) GetIsPerformanceAdvisorEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.IsPerformanceAdvisorEnabled) {
		return nil, false
	}

	return o.IsPerformanceAdvisorEnabled, true
}

// HasIsPerformanceAdvisorEnabled returns a boolean if a field has been set.
func (o *GroupSettings) HasIsPerformanceAdvisorEnabled() bool {
	if o != nil && !IsNil(o.IsPerformanceAdvisorEnabled) {
		return true
	}

	return false
}

// SetIsPerformanceAdvisorEnabled gets a reference to the given bool and assigns it to the IsPerformanceAdvisorEnabled field.
func (o *GroupSettings) SetIsPerformanceAdvisorEnabled(v bool) {
	o.IsPerformanceAdvisorEnabled = &v
}

// GetIsRealtimePerformancePanelEnabled returns the IsRealtimePerformancePanelEnabled field value if set, zero value otherwise
func (o *GroupSettings) GetIsRealtimePerformancePanelEnabled() bool {
	if o == nil || IsNil(o.IsRealtimePerformancePanelEnabled) {
		var ret bool
		return ret
	}
	return *o.IsRealtimePerformancePanelEnabled
}

// GetIsRealtimePerformancePanelEnabledOk returns a tuple with the IsRealtimePerformancePanelEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupSettings) GetIsRealtimePerformancePanelEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.IsRealtimePerformancePanelEnabled) {
		return nil, false
	}

	return o.IsRealtimePerformancePanelEnabled, true
}

// HasIsRealtimePerformancePanelEnabled returns a boolean if a field has been set.
func (o *GroupSettings) HasIsRealtimePerformancePanelEnabled() bool {
	if o != nil && !IsNil(o.IsRealtimePerformancePanelEnabled) {
		return true
	}

	return false
}

// SetIsRealtimePerformancePanelEnabled gets a reference to the given bool and assigns it to the IsRealtimePerformancePanelEnabled field.
func (o *GroupSettings) SetIsRealtimePerformancePanelEnabled(v bool) {
	o.IsRealtimePerformancePanelEnabled = &v
}

// GetIsSchemaAdvisorEnabled returns the IsSchemaAdvisorEnabled field value if set, zero value otherwise
func (o *GroupSettings) GetIsSchemaAdvisorEnabled() bool {
	if o == nil || IsNil(o.IsSchemaAdvisorEnabled) {
		var ret bool
		return ret
	}
	return *o.IsSchemaAdvisorEnabled
}

// GetIsSchemaAdvisorEnabledOk returns a tuple with the IsSchemaAdvisorEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupSettings) GetIsSchemaAdvisorEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.IsSchemaAdvisorEnabled) {
		return nil, false
	}

	return o.IsSchemaAdvisorEnabled, true
}

// HasIsSchemaAdvisorEnabled returns a boolean if a field has been set.
func (o *GroupSettings) HasIsSchemaAdvisorEnabled() bool {
	if o != nil && !IsNil(o.IsSchemaAdvisorEnabled) {
		return true
	}

	return false
}

// SetIsSchemaAdvisorEnabled gets a reference to the given bool and assigns it to the IsSchemaAdvisorEnabled field.
func (o *GroupSettings) SetIsSchemaAdvisorEnabled(v bool) {
	o.IsSchemaAdvisorEnabled = &v
}
