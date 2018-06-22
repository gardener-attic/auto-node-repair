// +build go1.9

// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package sql

import original "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"

type BackupShortTermRetentionPoliciesClient = original.BackupShortTermRetentionPoliciesClient
type CapabilitiesClient = original.CapabilitiesClient

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient
type DatabaseOperationsClient = original.DatabaseOperationsClient
type DatabasesClient = original.DatabasesClient
type DatabaseVulnerabilityAssessmentScansClient = original.DatabaseVulnerabilityAssessmentScansClient
type ElasticPoolOperationsClient = original.ElasticPoolOperationsClient
type ElasticPoolsClient = original.ElasticPoolsClient
type InstanceFailoverGroupsClient = original.InstanceFailoverGroupsClient
type CapabilityGroup = original.CapabilityGroup

const (
	SupportedEditions                CapabilityGroup = original.SupportedEditions
	SupportedElasticPoolEditions     CapabilityGroup = original.SupportedElasticPoolEditions
	SupportedManagedInstanceVersions CapabilityGroup = original.SupportedManagedInstanceVersions
)

type CapabilityStatus = original.CapabilityStatus

const (
	Available CapabilityStatus = original.Available
	Default   CapabilityStatus = original.Default
	Disabled  CapabilityStatus = original.Disabled
	Visible   CapabilityStatus = original.Visible
)

type CatalogCollationType = original.CatalogCollationType

const (
	DATABASEDEFAULT         CatalogCollationType = original.DATABASEDEFAULT
	SQLLatin1GeneralCP1CIAS CatalogCollationType = original.SQLLatin1GeneralCP1CIAS
)

type CreateMode = original.CreateMode

const (
	CreateModeCopy                           CreateMode = original.CreateModeCopy
	CreateModeDefault                        CreateMode = original.CreateModeDefault
	CreateModeOnlineSecondary                CreateMode = original.CreateModeOnlineSecondary
	CreateModePointInTimeRestore             CreateMode = original.CreateModePointInTimeRestore
	CreateModeRecovery                       CreateMode = original.CreateModeRecovery
	CreateModeRestore                        CreateMode = original.CreateModeRestore
	CreateModeRestoreExternalBackup          CreateMode = original.CreateModeRestoreExternalBackup
	CreateModeRestoreExternalBackupSecondary CreateMode = original.CreateModeRestoreExternalBackupSecondary
	CreateModeRestoreLongTermRetentionBackup CreateMode = original.CreateModeRestoreLongTermRetentionBackup
	CreateModeSecondary                      CreateMode = original.CreateModeSecondary
)

type DatabaseLicenseType = original.DatabaseLicenseType

const (
	BasePrice       DatabaseLicenseType = original.BasePrice
	LicenseIncluded DatabaseLicenseType = original.LicenseIncluded
)

type DatabaseReadScale = original.DatabaseReadScale

const (
	DatabaseReadScaleDisabled DatabaseReadScale = original.DatabaseReadScaleDisabled
	DatabaseReadScaleEnabled  DatabaseReadScale = original.DatabaseReadScaleEnabled
)

type DatabaseStatus = original.DatabaseStatus

const (
	AutoClosed       DatabaseStatus = original.AutoClosed
	Copying          DatabaseStatus = original.Copying
	Creating         DatabaseStatus = original.Creating
	EmergencyMode    DatabaseStatus = original.EmergencyMode
	Inaccessible     DatabaseStatus = original.Inaccessible
	Offline          DatabaseStatus = original.Offline
	OfflineSecondary DatabaseStatus = original.OfflineSecondary
	Online           DatabaseStatus = original.Online
	Paused           DatabaseStatus = original.Paused
	Pausing          DatabaseStatus = original.Pausing
	Recovering       DatabaseStatus = original.Recovering
	RecoveryPending  DatabaseStatus = original.RecoveryPending
	Restoring        DatabaseStatus = original.Restoring
	Resuming         DatabaseStatus = original.Resuming
	Scaling          DatabaseStatus = original.Scaling
	Shutdown         DatabaseStatus = original.Shutdown
	Standby          DatabaseStatus = original.Standby
	Suspect          DatabaseStatus = original.Suspect
)

type ElasticPoolLicenseType = original.ElasticPoolLicenseType

const (
	ElasticPoolLicenseTypeBasePrice       ElasticPoolLicenseType = original.ElasticPoolLicenseTypeBasePrice
	ElasticPoolLicenseTypeLicenseIncluded ElasticPoolLicenseType = original.ElasticPoolLicenseTypeLicenseIncluded
)

type ElasticPoolState = original.ElasticPoolState

const (
	ElasticPoolStateCreating ElasticPoolState = original.ElasticPoolStateCreating
	ElasticPoolStateDisabled ElasticPoolState = original.ElasticPoolStateDisabled
	ElasticPoolStateReady    ElasticPoolState = original.ElasticPoolStateReady
)

type InstanceFailoverGroupReplicationRole = original.InstanceFailoverGroupReplicationRole

const (
	Primary   InstanceFailoverGroupReplicationRole = original.Primary
	Secondary InstanceFailoverGroupReplicationRole = original.Secondary
)

type LogSizeUnit = original.LogSizeUnit

const (
	Gigabytes LogSizeUnit = original.Gigabytes
	Megabytes LogSizeUnit = original.Megabytes
	Percent   LogSizeUnit = original.Percent
	Petabytes LogSizeUnit = original.Petabytes
	Terabytes LogSizeUnit = original.Terabytes
)

type ManagementOperationState = original.ManagementOperationState

const (
	CancelInProgress ManagementOperationState = original.CancelInProgress
	Cancelled        ManagementOperationState = original.Cancelled
	Failed           ManagementOperationState = original.Failed
	InProgress       ManagementOperationState = original.InProgress
	Pending          ManagementOperationState = original.Pending
	Succeeded        ManagementOperationState = original.Succeeded
)

type MaxSizeUnit = original.MaxSizeUnit

const (
	MaxSizeUnitGigabytes MaxSizeUnit = original.MaxSizeUnitGigabytes
	MaxSizeUnitMegabytes MaxSizeUnit = original.MaxSizeUnitMegabytes
	MaxSizeUnitPetabytes MaxSizeUnit = original.MaxSizeUnitPetabytes
	MaxSizeUnitTerabytes MaxSizeUnit = original.MaxSizeUnitTerabytes
)

type PerformanceLevelUnit = original.PerformanceLevelUnit

const (
	DTU    PerformanceLevelUnit = original.DTU
	VCores PerformanceLevelUnit = original.VCores
)

type ReadOnlyEndpointFailoverPolicy = original.ReadOnlyEndpointFailoverPolicy

const (
	ReadOnlyEndpointFailoverPolicyDisabled ReadOnlyEndpointFailoverPolicy = original.ReadOnlyEndpointFailoverPolicyDisabled
	ReadOnlyEndpointFailoverPolicyEnabled  ReadOnlyEndpointFailoverPolicy = original.ReadOnlyEndpointFailoverPolicyEnabled
)

type ReadWriteEndpointFailoverPolicy = original.ReadWriteEndpointFailoverPolicy

const (
	Automatic ReadWriteEndpointFailoverPolicy = original.Automatic
	Manual    ReadWriteEndpointFailoverPolicy = original.Manual
)

type SampleName = original.SampleName

const (
	AdventureWorksLT       SampleName = original.AdventureWorksLT
	WideWorldImportersFull SampleName = original.WideWorldImportersFull
	WideWorldImportersStd  SampleName = original.WideWorldImportersStd
)

type VulnerabilityAssessmentScanState = original.VulnerabilityAssessmentScanState

const (
	VulnerabilityAssessmentScanStateFailed      VulnerabilityAssessmentScanState = original.VulnerabilityAssessmentScanStateFailed
	VulnerabilityAssessmentScanStateFailedToRun VulnerabilityAssessmentScanState = original.VulnerabilityAssessmentScanStateFailedToRun
	VulnerabilityAssessmentScanStateInProgress  VulnerabilityAssessmentScanState = original.VulnerabilityAssessmentScanStateInProgress
	VulnerabilityAssessmentScanStatePassed      VulnerabilityAssessmentScanState = original.VulnerabilityAssessmentScanStatePassed
)

type VulnerabilityAssessmentScanTriggerType = original.VulnerabilityAssessmentScanTriggerType

const (
	OnDemand  VulnerabilityAssessmentScanTriggerType = original.OnDemand
	Recurring VulnerabilityAssessmentScanTriggerType = original.Recurring
)

type BackupShortTermRetentionPoliciesCreateOrUpdateFuture = original.BackupShortTermRetentionPoliciesCreateOrUpdateFuture
type BackupShortTermRetentionPoliciesUpdateFuture = original.BackupShortTermRetentionPoliciesUpdateFuture
type BackupShortTermRetentionPolicy = original.BackupShortTermRetentionPolicy
type BackupShortTermRetentionPolicyListResult = original.BackupShortTermRetentionPolicyListResult
type BackupShortTermRetentionPolicyListResultIterator = original.BackupShortTermRetentionPolicyListResultIterator
type BackupShortTermRetentionPolicyListResultPage = original.BackupShortTermRetentionPolicyListResultPage
type BackupShortTermRetentionPolicyProperties = original.BackupShortTermRetentionPolicyProperties
type Database = original.Database
type DatabaseListResult = original.DatabaseListResult
type DatabaseListResultIterator = original.DatabaseListResultIterator
type DatabaseListResultPage = original.DatabaseListResultPage
type DatabaseOperation = original.DatabaseOperation
type DatabaseOperationListResult = original.DatabaseOperationListResult
type DatabaseOperationListResultIterator = original.DatabaseOperationListResultIterator
type DatabaseOperationListResultPage = original.DatabaseOperationListResultPage
type DatabaseOperationProperties = original.DatabaseOperationProperties
type DatabaseProperties = original.DatabaseProperties
type DatabasesCreateOrUpdateFuture = original.DatabasesCreateOrUpdateFuture
type DatabasesDeleteFuture = original.DatabasesDeleteFuture
type DatabasesPauseFuture = original.DatabasesPauseFuture
type DatabasesResumeFuture = original.DatabasesResumeFuture
type DatabasesUpdateFuture = original.DatabasesUpdateFuture
type DatabasesUpgradeDataWarehouseFuture = original.DatabasesUpgradeDataWarehouseFuture
type DatabaseUpdate = original.DatabaseUpdate
type DatabaseVulnerabilityAssessment = original.DatabaseVulnerabilityAssessment
type DatabaseVulnerabilityAssessmentProperties = original.DatabaseVulnerabilityAssessmentProperties
type DatabaseVulnerabilityAssessmentScanExportProperties = original.DatabaseVulnerabilityAssessmentScanExportProperties
type DatabaseVulnerabilityAssessmentScansExport = original.DatabaseVulnerabilityAssessmentScansExport
type DatabaseVulnerabilityAssessmentScansInitiateScanFuture = original.DatabaseVulnerabilityAssessmentScansInitiateScanFuture
type EditionCapability = original.EditionCapability
type ElasticPool = original.ElasticPool
type ElasticPoolEditionCapability = original.ElasticPoolEditionCapability
type ElasticPoolListResult = original.ElasticPoolListResult
type ElasticPoolListResultIterator = original.ElasticPoolListResultIterator
type ElasticPoolListResultPage = original.ElasticPoolListResultPage
type ElasticPoolOperation = original.ElasticPoolOperation
type ElasticPoolOperationListResult = original.ElasticPoolOperationListResult
type ElasticPoolOperationListResultIterator = original.ElasticPoolOperationListResultIterator
type ElasticPoolOperationListResultPage = original.ElasticPoolOperationListResultPage
type ElasticPoolOperationProperties = original.ElasticPoolOperationProperties
type ElasticPoolPerDatabaseMaxPerformanceLevelCapability = original.ElasticPoolPerDatabaseMaxPerformanceLevelCapability
type ElasticPoolPerDatabaseMinPerformanceLevelCapability = original.ElasticPoolPerDatabaseMinPerformanceLevelCapability
type ElasticPoolPerDatabaseSettings = original.ElasticPoolPerDatabaseSettings
type ElasticPoolPerformanceLevelCapability = original.ElasticPoolPerformanceLevelCapability
type ElasticPoolProperties = original.ElasticPoolProperties
type ElasticPoolsCreateOrUpdateFuture = original.ElasticPoolsCreateOrUpdateFuture
type ElasticPoolsDeleteFuture = original.ElasticPoolsDeleteFuture
type ElasticPoolsUpdateFuture = original.ElasticPoolsUpdateFuture
type ElasticPoolUpdate = original.ElasticPoolUpdate
type ElasticPoolUpdateProperties = original.ElasticPoolUpdateProperties
type InstanceFailoverGroup = original.InstanceFailoverGroup
type InstanceFailoverGroupListResult = original.InstanceFailoverGroupListResult
type InstanceFailoverGroupListResultIterator = original.InstanceFailoverGroupListResultIterator
type InstanceFailoverGroupListResultPage = original.InstanceFailoverGroupListResultPage
type InstanceFailoverGroupProperties = original.InstanceFailoverGroupProperties
type InstanceFailoverGroupReadOnlyEndpoint = original.InstanceFailoverGroupReadOnlyEndpoint
type InstanceFailoverGroupReadWriteEndpoint = original.InstanceFailoverGroupReadWriteEndpoint
type InstanceFailoverGroupsCreateOrUpdateFuture = original.InstanceFailoverGroupsCreateOrUpdateFuture
type InstanceFailoverGroupsDeleteFuture = original.InstanceFailoverGroupsDeleteFuture
type InstanceFailoverGroupsFailoverFuture = original.InstanceFailoverGroupsFailoverFuture
type InstanceFailoverGroupsForceFailoverAllowDataLossFuture = original.InstanceFailoverGroupsForceFailoverAllowDataLossFuture
type LicenseTypeCapability = original.LicenseTypeCapability
type LocationCapabilities = original.LocationCapabilities
type LogSizeCapability = original.LogSizeCapability
type ManagedInstanceEditionCapability = original.ManagedInstanceEditionCapability
type ManagedInstanceFamilyCapability = original.ManagedInstanceFamilyCapability
type ManagedInstancePairInfo = original.ManagedInstancePairInfo
type ManagedInstanceVcoresCapability = original.ManagedInstanceVcoresCapability
type ManagedInstanceVersionCapability = original.ManagedInstanceVersionCapability
type MaxSizeCapability = original.MaxSizeCapability
type MaxSizeRangeCapability = original.MaxSizeRangeCapability
type PartnerRegionInfo = original.PartnerRegionInfo
type PerformanceLevelCapability = original.PerformanceLevelCapability
type ProxyResource = original.ProxyResource
type Resource = original.Resource
type ResourceMoveDefinition = original.ResourceMoveDefinition
type ServerVersionCapability = original.ServerVersionCapability
type ServiceObjectiveCapability = original.ServiceObjectiveCapability
type Sku = original.Sku
type TrackedResource = original.TrackedResource
type VulnerabilityAssessmentRecurringScansProperties = original.VulnerabilityAssessmentRecurringScansProperties
type VulnerabilityAssessmentScanError = original.VulnerabilityAssessmentScanError
type VulnerabilityAssessmentScanRecord = original.VulnerabilityAssessmentScanRecord
type VulnerabilityAssessmentScanRecordListResult = original.VulnerabilityAssessmentScanRecordListResult
type VulnerabilityAssessmentScanRecordListResultIterator = original.VulnerabilityAssessmentScanRecordListResultIterator
type VulnerabilityAssessmentScanRecordListResultPage = original.VulnerabilityAssessmentScanRecordListResultPage
type VulnerabilityAssessmentScanRecordProperties = original.VulnerabilityAssessmentScanRecordProperties

func NewBackupShortTermRetentionPoliciesClient(subscriptionID string) BackupShortTermRetentionPoliciesClient {
	return original.NewBackupShortTermRetentionPoliciesClient(subscriptionID)
}
func NewBackupShortTermRetentionPoliciesClientWithBaseURI(baseURI string, subscriptionID string) BackupShortTermRetentionPoliciesClient {
	return original.NewBackupShortTermRetentionPoliciesClientWithBaseURI(baseURI, subscriptionID)
}
func NewCapabilitiesClient(subscriptionID string) CapabilitiesClient {
	return original.NewCapabilitiesClient(subscriptionID)
}
func NewCapabilitiesClientWithBaseURI(baseURI string, subscriptionID string) CapabilitiesClient {
	return original.NewCapabilitiesClientWithBaseURI(baseURI, subscriptionID)
}
func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}
func NewDatabaseOperationsClient(subscriptionID string) DatabaseOperationsClient {
	return original.NewDatabaseOperationsClient(subscriptionID)
}
func NewDatabaseOperationsClientWithBaseURI(baseURI string, subscriptionID string) DatabaseOperationsClient {
	return original.NewDatabaseOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewDatabasesClient(subscriptionID string) DatabasesClient {
	return original.NewDatabasesClient(subscriptionID)
}
func NewDatabasesClientWithBaseURI(baseURI string, subscriptionID string) DatabasesClient {
	return original.NewDatabasesClientWithBaseURI(baseURI, subscriptionID)
}
func NewDatabaseVulnerabilityAssessmentScansClient(subscriptionID string) DatabaseVulnerabilityAssessmentScansClient {
	return original.NewDatabaseVulnerabilityAssessmentScansClient(subscriptionID)
}
func NewDatabaseVulnerabilityAssessmentScansClientWithBaseURI(baseURI string, subscriptionID string) DatabaseVulnerabilityAssessmentScansClient {
	return original.NewDatabaseVulnerabilityAssessmentScansClientWithBaseURI(baseURI, subscriptionID)
}
func NewElasticPoolOperationsClient(subscriptionID string) ElasticPoolOperationsClient {
	return original.NewElasticPoolOperationsClient(subscriptionID)
}
func NewElasticPoolOperationsClientWithBaseURI(baseURI string, subscriptionID string) ElasticPoolOperationsClient {
	return original.NewElasticPoolOperationsClientWithBaseURI(baseURI, subscriptionID)
}
func NewElasticPoolsClient(subscriptionID string) ElasticPoolsClient {
	return original.NewElasticPoolsClient(subscriptionID)
}
func NewElasticPoolsClientWithBaseURI(baseURI string, subscriptionID string) ElasticPoolsClient {
	return original.NewElasticPoolsClientWithBaseURI(baseURI, subscriptionID)
}
func NewInstanceFailoverGroupsClient(subscriptionID string) InstanceFailoverGroupsClient {
	return original.NewInstanceFailoverGroupsClient(subscriptionID)
}
func NewInstanceFailoverGroupsClientWithBaseURI(baseURI string, subscriptionID string) InstanceFailoverGroupsClient {
	return original.NewInstanceFailoverGroupsClientWithBaseURI(baseURI, subscriptionID)
}
func PossibleCapabilityGroupValues() []CapabilityGroup {
	return original.PossibleCapabilityGroupValues()
}
func PossibleCapabilityStatusValues() []CapabilityStatus {
	return original.PossibleCapabilityStatusValues()
}
func PossibleCatalogCollationTypeValues() []CatalogCollationType {
	return original.PossibleCatalogCollationTypeValues()
}
func PossibleCreateModeValues() []CreateMode {
	return original.PossibleCreateModeValues()
}
func PossibleDatabaseLicenseTypeValues() []DatabaseLicenseType {
	return original.PossibleDatabaseLicenseTypeValues()
}
func PossibleDatabaseReadScaleValues() []DatabaseReadScale {
	return original.PossibleDatabaseReadScaleValues()
}
func PossibleDatabaseStatusValues() []DatabaseStatus {
	return original.PossibleDatabaseStatusValues()
}
func PossibleElasticPoolLicenseTypeValues() []ElasticPoolLicenseType {
	return original.PossibleElasticPoolLicenseTypeValues()
}
func PossibleElasticPoolStateValues() []ElasticPoolState {
	return original.PossibleElasticPoolStateValues()
}
func PossibleInstanceFailoverGroupReplicationRoleValues() []InstanceFailoverGroupReplicationRole {
	return original.PossibleInstanceFailoverGroupReplicationRoleValues()
}
func PossibleLogSizeUnitValues() []LogSizeUnit {
	return original.PossibleLogSizeUnitValues()
}
func PossibleManagementOperationStateValues() []ManagementOperationState {
	return original.PossibleManagementOperationStateValues()
}
func PossibleMaxSizeUnitValues() []MaxSizeUnit {
	return original.PossibleMaxSizeUnitValues()
}
func PossiblePerformanceLevelUnitValues() []PerformanceLevelUnit {
	return original.PossiblePerformanceLevelUnitValues()
}
func PossibleReadOnlyEndpointFailoverPolicyValues() []ReadOnlyEndpointFailoverPolicy {
	return original.PossibleReadOnlyEndpointFailoverPolicyValues()
}
func PossibleReadWriteEndpointFailoverPolicyValues() []ReadWriteEndpointFailoverPolicy {
	return original.PossibleReadWriteEndpointFailoverPolicyValues()
}
func PossibleSampleNameValues() []SampleName {
	return original.PossibleSampleNameValues()
}
func PossibleVulnerabilityAssessmentScanStateValues() []VulnerabilityAssessmentScanState {
	return original.PossibleVulnerabilityAssessmentScanStateValues()
}
func PossibleVulnerabilityAssessmentScanTriggerTypeValues() []VulnerabilityAssessmentScanTriggerType {
	return original.PossibleVulnerabilityAssessmentScanTriggerTypeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
