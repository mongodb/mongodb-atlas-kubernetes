# API Reference

Packages:

- [atlas.mongodb.com/v1](#atlasmongodbcomv1)

# atlas.mongodb.com/v1

Resource Types:

- [AtlasBackupCompliancePolicy](#atlasbackupcompliancepolicy)

- [AtlasBackupPolicy](#atlasbackuppolicy)

- [AtlasBackupSchedule](#atlasbackupschedule)

- [AtlasCustomRole](#atlascustomrole)

- [AtlasDatabaseUser](#atlasdatabaseuser)

- [AtlasDataFederation](#atlasdatafederation)

- [AtlasDeployment](#atlasdeployment)

- [AtlasFederatedAuth](#atlasfederatedauth)

- [AtlasIPAccessList](#atlasipaccesslist)

- [AtlasNetworkContainer](#atlasnetworkcontainer)

- [AtlasNetworkPeering](#atlasnetworkpeering)

- [AtlasOrgSettings](#atlasorgsettings)

- [AtlasPrivateEndpoint](#atlasprivateendpoint)

- [AtlasProject](#atlasproject)

- [AtlasSearchIndexConfig](#atlassearchindexconfig)

- [AtlasStreamConnection](#atlasstreamconnection)

- [AtlasStreamInstance](#atlasstreaminstance)

- [AtlasTeam](#atlasteam)

- [AtlasThirdPartyIntegration](#atlasthirdpartyintegration)




## AtlasBackupCompliancePolicy
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






The AtlasBackupCompliancePolicy is a configuration that enforces specific backup and retention requirements

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasBackupCompliancePolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasbackupcompliancepolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasBackupCompliancePolicySpec is the specification of the desired configuration of backup compliance policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasbackupcompliancepolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          BackupCompliancePolicyStatus defines the observed state of AtlasBackupCompliancePolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasBackupCompliancePolicy.spec
<sup><sup>[↩ Parent](#atlasbackupcompliancepolicy)</sup></sup>



AtlasBackupCompliancePolicySpec is the specification of the desired configuration of backup compliance policy

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>authorizedEmail</b></td>
        <td>string</td>
        <td>
          Email address of the user who authorized to update the Backup Compliance Policy settings.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>authorizedUserFirstName</b></td>
        <td>string</td>
        <td>
          First name of the user who authorized to updated the Backup Compliance Policy settings.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>authorizedUserLastName</b></td>
        <td>string</td>
        <td>
          Last name of the user who authorized to updated the Backup Compliance Policy settings.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>copyProtectionEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to prevent cluster users from deleting backups copied to other regions, even if those additional snapshot regions are removed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>encryptionAtRestEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether Encryption at Rest using Customer Key Management is required for all clusters with a Backup Compliance Policy.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasbackupcompliancepolicyspecondemandpolicy">onDemandPolicy</a></b></td>
        <td>object</td>
        <td>
          Specifications for on-demand policy.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>overwriteBackupPolicies</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to overwrite non-complying backup policies with the new data protection settings or not.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>pointInTimeEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether the cluster uses Continuous Cloud Backups with a Backup Compliance Policy.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>restoreWindowDays</b></td>
        <td>integer</td>
        <td>
          Number of previous days that you can restore back to with Continuous Cloud Backup with a Backup Compliance Policy. This parameter applies only to Continuous Cloud Backups with a Backup Compliance Policy.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasbackupcompliancepolicyspecscheduledpolicyitemsindex">scheduledPolicyItems</a></b></td>
        <td>[]object</td>
        <td>
          List that contains the specifications for one scheduled policy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasBackupCompliancePolicy.spec.onDemandPolicy
<sup><sup>[↩ Parent](#atlasbackupcompliancepolicyspec)</sup></sup>



Specifications for on-demand policy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>retentionUnit</b></td>
        <td>enum</td>
        <td>
          Scope of the backup policy item: days, weeks, or months.<br/>
          <br/>
            <i>Enum</i>: days, weeks, months<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>retentionValue</b></td>
        <td>integer</td>
        <td>
          Value to associate with RetentionUnit.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasBackupCompliancePolicy.spec.scheduledPolicyItems[index]
<sup><sup>[↩ Parent](#atlasbackupcompliancepolicyspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>frequencyInterval</b></td>
        <td>integer</td>
        <td>
          Desired frequency of the new backup policy item specified by FrequencyType. A value of 1 specifies the first instance of the corresponding FrequencyType.
The only accepted value you can set for frequency interval with NVMe clusters is 12.<br/>
          <br/>
            <i>Enum</i>: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 40<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>frequencyType</b></td>
        <td>enum</td>
        <td>
          Frequency associated with the backup policy item. You cannot specify multiple hourly and daily backup policy items.<br/>
          <br/>
            <i>Enum</i>: hourly, daily, weekly, monthly, yearly<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>retentionUnit</b></td>
        <td>enum</td>
        <td>
          Unit of time in which MongoDB Atlas measures snapshot retention.<br/>
          <br/>
            <i>Enum</i>: days, weeks, months, years<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>retentionValue</b></td>
        <td>integer</td>
        <td>
          Duration in days, weeks, months, or years that MongoDB Cloud retains the snapshot.
For less frequent policy items, MongoDB Cloud requires that you specify a value greater than or equal to the value specified for more frequent policy items.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasBackupCompliancePolicy.status
<sup><sup>[↩ Parent](#atlasbackupcompliancepolicy)</sup></sup>



BackupCompliancePolicyStatus defines the observed state of AtlasBackupCompliancePolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasbackupcompliancepolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasBackupCompliancePolicy.status.conditions[index]
<sup><sup>[↩ Parent](#atlasbackupcompliancepolicystatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasBackupPolicy
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasBackupPolicy is the Schema for the atlasbackuppolicies API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasBackupPolicy</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasbackuppolicyspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasBackupPolicySpec defines the desired state of AtlasBackupPolicy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasbackuppolicystatus">status</a></b></td>
        <td>object</td>
        <td>
          BackupPolicyStatus defines the observed state of AtlasBackupPolicy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasBackupPolicy.spec
<sup><sup>[↩ Parent](#atlasbackuppolicy)</sup></sup>



AtlasBackupPolicySpec defines the desired state of AtlasBackupPolicy

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasbackuppolicyspecitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          A list of BackupPolicy items.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasBackupPolicy.spec.items[index]
<sup><sup>[↩ Parent](#atlasbackuppolicyspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>frequencyInterval</b></td>
        <td>integer</td>
        <td>
          Desired frequency of the new backup policy item specified by FrequencyType. A value of 1 specifies the first instance of the corresponding FrequencyType.
The only accepted value you can set for frequency interval with NVMe clusters is 12.<br/>
          <br/>
            <i>Enum</i>: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 40<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>frequencyType</b></td>
        <td>enum</td>
        <td>
          Frequency associated with the backup policy item. You cannot specify multiple hourly and daily backup policy items.<br/>
          <br/>
            <i>Enum</i>: hourly, daily, weekly, monthly, yearly<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>retentionUnit</b></td>
        <td>enum</td>
        <td>
          Unit of time in which MongoDB Atlas measures snapshot retention.<br/>
          <br/>
            <i>Enum</i>: days, weeks, months, years<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>retentionValue</b></td>
        <td>integer</td>
        <td>
          Duration in days, weeks, months, or years that MongoDB Cloud retains the snapshot.
For less frequent policy items, MongoDB Cloud requires that you specify a value greater than or equal to the value specified for more frequent policy items.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasBackupPolicy.status
<sup><sup>[↩ Parent](#atlasbackuppolicy)</sup></sup>



BackupPolicyStatus defines the observed state of AtlasBackupPolicy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasbackuppolicystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>backupScheduleIDs</b></td>
        <td>[]string</td>
        <td>
          DeploymentID of the deployment using the backup policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasBackupPolicy.status.conditions[index]
<sup><sup>[↩ Parent](#atlasbackuppolicystatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasBackupSchedule
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasBackupSchedule is the Schema for the atlasbackupschedules API.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasBackupSchedule</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasbackupschedulespec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasBackupScheduleSpec defines the desired state of AtlasBackupSchedule.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasbackupschedulestatus">status</a></b></td>
        <td>object</td>
        <td>
          BackupScheduleStatus defines the observed state of AtlasBackupSchedule.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasBackupSchedule.spec
<sup><sup>[↩ Parent](#atlasbackupschedule)</sup></sup>



AtlasBackupScheduleSpec defines the desired state of AtlasBackupSchedule.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasbackupschedulespecpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          A reference (name & namespace) for backup policy in the desired updated backup policy.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>autoExportEnabled</b></td>
        <td>boolean</td>
        <td>
          Specify true to enable automatic export of cloud backup snapshots to the AWS bucket. You must also define the export policy using export. If omitted, defaults to false.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasbackupschedulespeccopysettingsindex">copySettings</a></b></td>
        <td>[]object</td>
        <td>
          Copy backups to other regions for increased resiliency and faster restores.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasbackupschedulespecexport">export</a></b></td>
        <td>object</td>
        <td>
          Export policy for automatically exporting cloud backup snapshots to AWS bucket.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>referenceHourOfDay</b></td>
        <td>integer</td>
        <td>
          UTC Hour of day between 0 and 23, inclusive, representing which hour of the day that Atlas takes snapshots for backup policy items<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Minimum</i>: 0<br/>
            <i>Maximum</i>: 23<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>referenceMinuteOfHour</b></td>
        <td>integer</td>
        <td>
          UTC Minutes after ReferenceHourOfDay that Atlas takes snapshots for backup policy items. Must be between 0 and 59, inclusive.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Minimum</i>: 0<br/>
            <i>Maximum</i>: 59<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>restoreWindowDays</b></td>
        <td>integer</td>
        <td>
          Number of days back in time you can restore to with Continuous Cloud Backup accuracy. Must be a positive, non-zero integer. Applies to continuous cloud backups only.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Default</i>: 1<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>updateSnapshots</b></td>
        <td>boolean</td>
        <td>
          Specify true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>useOrgAndGroupNamesInExportPrefix</b></td>
        <td>boolean</td>
        <td>
          Specify true to use organization and project names instead of organization and project UUIDs in the path for the metadata files that Atlas uploads to your S3 bucket after it finishes exporting the snapshots<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasBackupSchedule.spec.policy
<sup><sup>[↩ Parent](#atlasbackupschedulespec)</sup></sup>



A reference (name & namespace) for backup policy in the desired updated backup policy.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasBackupSchedule.spec.copySettings[index]
<sup><sup>[↩ Parent](#atlasbackupschedulespec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>cloudProvider</b></td>
        <td>enum</td>
        <td>
          Identifies the cloud provider that stores the snapshot copy.<br/>
          <br/>
            <i>Enum</i>: AWS, GCP, AZURE<br/>
            <i>Default</i>: AWS<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>frequencies</b></td>
        <td>[]string</td>
        <td>
          List that describes which types of snapshots to copy.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>regionName</b></td>
        <td>string</td>
        <td>
          Target region to copy snapshots belonging to replicationSpecId to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>shouldCopyOplogs</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to copy the oplogs to the target region.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasBackupSchedule.spec.export
<sup><sup>[↩ Parent](#atlasbackupschedulespec)</sup></sup>



Export policy for automatically exporting cloud backup snapshots to AWS bucket.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>exportBucketId</b></td>
        <td>string</td>
        <td>
          Unique Atlas identifier of the AWS bucket which was granted access to export backup snapshot.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>frequencyType</b></td>
        <td>enum</td>
        <td>
          Human-readable label that indicates the rate at which the export policy item occurs.<br/>
          <br/>
            <i>Enum</i>: monthly<br/>
            <i>Default</i>: monthly<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasBackupSchedule.status
<sup><sup>[↩ Parent](#atlasbackupschedule)</sup></sup>



BackupScheduleStatus defines the observed state of AtlasBackupSchedule.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasbackupschedulestatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>deploymentID</b></td>
        <td>[]string</td>
        <td>
          List of the human-readable names of all deployments utilizing this backup schedule.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasBackupSchedule.status.conditions[index]
<sup><sup>[↩ Parent](#atlasbackupschedulestatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasCustomRole
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasCustomRole is the Schema for the AtlasCustomRole API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasCustomRole</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlascustomrolespec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasCustomRoleSpec defines the desired state of CustomRole in Atlas.<br/>
          <br/>
            <i>Validations</i>:<li>(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef)): must define only one project reference through externalProjectRef or projectRef</li><li>(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef): must define a local connection secret when referencing an external project</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlascustomrolestatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasCustomRoleStatus is a status for the AtlasCustomRole Custom resource.
Not the one included in the AtlasProject<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasCustomRole.spec
<sup><sup>[↩ Parent](#atlascustomrole)</sup></sup>



AtlasCustomRoleSpec defines the desired state of CustomRole in Atlas.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlascustomrolespecrole">role</a></b></td>
        <td>object</td>
        <td>
          Role represents a Custom Role in Atlas.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlascustomrolespecconnectionsecret">connectionSecret</a></b></td>
        <td>object</td>
        <td>
          Name of the secret containing Atlas API private and public keys.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlascustomrolespecexternalprojectref">externalProjectRef</a></b></td>
        <td>object</td>
        <td>
          externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlascustomrolespecprojectref">projectRef</a></b></td>
        <td>object</td>
        <td>
          projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasCustomRole.spec.role
<sup><sup>[↩ Parent](#atlascustomrolespec)</sup></sup>



Role represents a Custom Role in Atlas.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the role. This name must be unique for this custom role in this project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlascustomrolespecroleactionsindex">actions</a></b></td>
        <td>[]object</td>
        <td>
          List of the individual privilege actions that the role grants.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlascustomrolespecroleinheritedrolesindex">inheritedRoles</a></b></td>
        <td>[]object</td>
        <td>
          List of the built-in roles that this custom role inherits.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasCustomRole.spec.role.actions[index]
<sup><sup>[↩ Parent](#atlascustomrolespecrole)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the privilege action.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlascustomrolespecroleactionsindexresourcesindex">resources</a></b></td>
        <td>[]object</td>
        <td>
          List of resources on which you grant the action.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasCustomRole.spec.role.actions[index].resources[index]
<sup><sup>[↩ Parent](#atlascustomrolespecroleactionsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>cluster</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to grant the action on the cluster resource. If true, MongoDB Cloud ignores Database and Collection parameters.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>collection</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the collection on which you grant the action to one MongoDB user.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the database on which you grant the action to one MongoDB user.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasCustomRole.spec.role.inheritedRoles[index]
<sup><sup>[↩ Parent](#atlascustomrolespecrole)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the database on which someone grants the action to one MongoDB user.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the role inherited.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasCustomRole.spec.connectionSecret
<sup><sup>[↩ Parent](#atlascustomrolespec)</sup></sup>



Name of the secret containing Atlas API private and public keys.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasCustomRole.spec.externalProjectRef
<sup><sup>[↩ Parent](#atlascustomrolespec)</sup></sup>



externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID is the Atlas project ID.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasCustomRole.spec.projectRef
<sup><sup>[↩ Parent](#atlascustomrolespec)</sup></sup>



projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasCustomRole.status
<sup><sup>[↩ Parent](#atlascustomrole)</sup></sup>



AtlasCustomRoleStatus is a status for the AtlasCustomRole Custom resource.
Not the one included in the AtlasProject

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlascustomrolestatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasCustomRole.status.conditions[index]
<sup><sup>[↩ Parent](#atlascustomrolestatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasDatabaseUser
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasDatabaseUser is the Schema for the Atlas Database User API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasDatabaseUser</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdatabaseuserspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasDatabaseUserSpec defines the desired state of Database User in Atlas<br/>
          <br/>
            <i>Validations</i>:<li>(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef)): must define only one project reference through externalProjectRef or projectRef</li><li>(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef): must define a local connection secret when referencing an external project</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatabaseuserstatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasDatabaseUserStatus defines the observed state of AtlasProject<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDatabaseUser.spec
<sup><sup>[↩ Parent](#atlasdatabaseuser)</sup></sup>



AtlasDatabaseUserSpec defines the desired state of Database User in Atlas

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdatabaseuserspecrolesindex">roles</a></b></td>
        <td>[]object</td>
        <td>
          Roles is an array of this user's roles and the databases / collections on which the roles apply. A role allows
the user to perform particular actions on the specified database.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>username</b></td>
        <td>string</td>
        <td>
          Username is a username for authenticating to MongoDB
Human-readable label that represents the user that authenticates to MongoDB. The format of this label depends on the method of authentication:
In case of AWS IAM: the value should be AWS ARN for the IAM User/Role;
In case of OIDC Workload or Workforce: the value should be the Atlas OIDC IdP ID, followed by a '/', followed by the IdP group name;
In case of Plain text auth: the value can be anything.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>awsIamType</b></td>
        <td>enum</td>
        <td>
          Human-readable label that indicates whether the new database user authenticates with Amazon Web Services (AWS).
Identity and Access Management (IAM) credentials associated with the user or the user's role<br/>
          <br/>
            <i>Enum</i>: NONE, USER, ROLE<br/>
            <i>Default</i>: NONE<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatabaseuserspecconnectionsecret">connectionSecret</a></b></td>
        <td>object</td>
        <td>
          Name of the secret containing Atlas API private and public keys.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>databaseName</b></td>
        <td>string</td>
        <td>
          DatabaseName is a Database against which Atlas authenticates the user.
If the user authenticates with AWS IAM, x.509, LDAP, or OIDC Workload this value should be '$external'.
If the user authenticates with SCRAM-SHA or OIDC Workforce, this value should be 'admin'.
Default value is 'admin'.<br/>
          <br/>
            <i>Default</i>: admin<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deleteAfterDate</b></td>
        <td>string</td>
        <td>
          DeleteAfterDate is a timestamp in ISO 8601 date and time format in UTC after which Atlas deletes the user.
The specified date must be in the future and within one week.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          Description of this database user. Maximum 100 characters.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatabaseuserspecexternalprojectref">externalProjectRef</a></b></td>
        <td>object</td>
        <td>
          externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatabaseuserspeclabelsindex">labels</a></b></td>
        <td>[]object</td>
        <td>
          Labels is an array containing key-value pairs that tag and categorize the database user.
Each key and value has a maximum length of 255 characters.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>oidcAuthType</b></td>
        <td>enum</td>
        <td>
          Human-readable label that indicates whether the new database Username with OIDC federated authentication.
To create a federated authentication group (Workforce), specify the value of IDP_GROUP in this field.
To create a federated authentication user (Workload), specify the value of USER in this field.<br/>
          <br/>
            <i>Enum</i>: NONE, IDP_GROUP, USER<br/>
            <i>Default</i>: NONE<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatabaseuserspecpasswordsecretref">passwordSecretRef</a></b></td>
        <td>object</td>
        <td>
          PasswordSecret is a reference to the Secret keeping the user password.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatabaseuserspecprojectref">projectRef</a></b></td>
        <td>object</td>
        <td>
          projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatabaseuserspecscopesindex">scopes</a></b></td>
        <td>[]object</td>
        <td>
          Scopes is an array of clusters and Atlas Data Lakes that this user has access to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>x509Type</b></td>
        <td>enum</td>
        <td>
          X509Type is X.509 method by which the database authenticates the provided username.<br/>
          <br/>
            <i>Enum</i>: NONE, MANAGED, CUSTOMER<br/>
            <i>Default</i>: NONE<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDatabaseUser.spec.roles[index]
<sup><sup>[↩ Parent](#atlasdatabaseuserspec)</sup></sup>



RoleSpec allows the user to perform particular actions on the specified database.
A role on the admin database can include privileges that apply to the other databases as well.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>databaseName</b></td>
        <td>string</td>
        <td>
          DatabaseName is a database on which the user has the specified role. A role on the admin database can include
privileges that apply to the other databases.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>roleName</b></td>
        <td>string</td>
        <td>
          RoleName is a name of the role. This value can either be a built-in role or a custom role.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>collectionName</b></td>
        <td>string</td>
        <td>
          CollectionName is a collection for which the role applies.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDatabaseUser.spec.connectionSecret
<sup><sup>[↩ Parent](#atlasdatabaseuserspec)</sup></sup>



Name of the secret containing Atlas API private and public keys.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDatabaseUser.spec.externalProjectRef
<sup><sup>[↩ Parent](#atlasdatabaseuserspec)</sup></sup>



externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID is the Atlas project ID.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDatabaseUser.spec.labels[index]
<sup><sup>[↩ Parent](#atlasdatabaseuserspec)</sup></sup>



LabelSpec contains key-value pairs that tag and categorize the Cluster/DBUser

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Key applied to tag and categorize this component.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Value set to the Key applied to tag and categorize this component.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDatabaseUser.spec.passwordSecretRef
<sup><sup>[↩ Parent](#atlasdatabaseuserspec)</sup></sup>



PasswordSecret is a reference to the Secret keeping the user password.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDatabaseUser.spec.projectRef
<sup><sup>[↩ Parent](#atlasdatabaseuserspec)</sup></sup>



projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDatabaseUser.spec.scopes[index]
<sup><sup>[↩ Parent](#atlasdatabaseuserspec)</sup></sup>



ScopeSpec if present a database user only have access to the indicated resource (Cluster or Atlas Data Lake)
if none is given then it has access to all.
It's highly recommended to restrict the access of the database users only to a limited set of resources.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is a name of the cluster or Atlas Data Lake that the user has access to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Type is a type of resource that the user has access to.<br/>
          <br/>
            <i>Enum</i>: CLUSTER, DATA_LAKE<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDatabaseUser.status
<sup><sup>[↩ Parent](#atlasdatabaseuser)</sup></sup>



AtlasDatabaseUserStatus defines the observed state of AtlasProject

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdatabaseuserstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          UserName is the current name of database user.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>passwordVersion</b></td>
        <td>string</td>
        <td>
          PasswordVersion is the 'ResourceVersion' of the password Secret that the Atlas Operator is aware of<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDatabaseUser.status.conditions[index]
<sup><sup>[↩ Parent](#atlasdatabaseuserstatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasDataFederation
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasDataFederation is the Schema for the Atlas Data Federation API.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasDataFederation</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdatafederationspec">spec</a></b></td>
        <td>object</td>
        <td>
          DataFederationSpec defines the desired state of AtlasDataFederation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatafederationstatus">status</a></b></td>
        <td>object</td>
        <td>
          DataFederationStatus defines the observed state of AtlasDataFederation.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec
<sup><sup>[↩ Parent](#atlasdatafederation)</sup></sup>



DataFederationSpec defines the desired state of AtlasDataFederation.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the Federated Database Instance.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdatafederationspecprojectref">projectRef</a></b></td>
        <td>object</td>
        <td>
          Project is a reference to AtlasProject resource the deployment belongs to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdatafederationspeccloudproviderconfig">cloudProviderConfig</a></b></td>
        <td>object</td>
        <td>
          Configuration for the cloud provider where this Federated Database Instance is hosted.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatafederationspecdataprocessregion">dataProcessRegion</a></b></td>
        <td>object</td>
        <td>
          Information about the cloud provider region to which the Federated Database Instance routes client connections.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatafederationspecprivateendpointsindex">privateEndpoints</a></b></td>
        <td>[]object</td>
        <td>
          Private endpoint for Federated Database Instances and Online Archives to add to the specified project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatafederationspecstorage">storage</a></b></td>
        <td>object</td>
        <td>
          Configuration information for each data store and its mapping to MongoDB Atlas databases.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec.projectRef
<sup><sup>[↩ Parent](#atlasdatafederationspec)</sup></sup>



Project is a reference to AtlasProject resource the deployment belongs to.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec.cloudProviderConfig
<sup><sup>[↩ Parent](#atlasdatafederationspec)</sup></sup>



Configuration for the cloud provider where this Federated Database Instance is hosted.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdatafederationspeccloudproviderconfigaws">aws</a></b></td>
        <td>object</td>
        <td>
          Configuration for running Data Federation in AWS.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec.cloudProviderConfig.aws
<sup><sup>[↩ Parent](#atlasdatafederationspeccloudproviderconfig)</sup></sup>



Configuration for running Data Federation in AWS.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>roleId</b></td>
        <td>string</td>
        <td>
          Unique identifier of the role that the data lake can use to access the data stores.Required if specifying cloudProviderConfig.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>testS3Bucket</b></td>
        <td>string</td>
        <td>
          Name of the S3 data bucket that the provided role ID is authorized to access.Required if specifying cloudProviderConfig.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec.dataProcessRegion
<sup><sup>[↩ Parent](#atlasdatafederationspec)</sup></sup>



Information about the cloud provider region to which the Federated Database Instance routes client connections.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>cloudProvider</b></td>
        <td>enum</td>
        <td>
          Name of the cloud service that hosts the Federated Database Instance's infrastructure.<br/>
          <br/>
            <i>Enum</i>: AWS<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>enum</td>
        <td>
          Name of the region to which the data lake routes client connections.<br/>
          <br/>
            <i>Enum</i>: SYDNEY_AUS, MUMBAI_IND, FRANKFURT_DEU, DUBLIN_IRL, LONDON_GBR, VIRGINIA_USA, OREGON_USA, SAOPAULO_BRA, SINGAPORE_SGP<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec.privateEndpoints[index]
<sup><sup>[↩ Parent](#atlasdatafederationspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>endpointId</b></td>
        <td>string</td>
        <td>
          Unique 22-character alphanumeric string that identifies the private endpoint.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>provider</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the cloud service provider. Atlas Data Lake supports Amazon Web Services only.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the resource type associated with this private endpoint.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec.storage
<sup><sup>[↩ Parent](#atlasdatafederationspec)</sup></sup>



Configuration information for each data store and its mapping to MongoDB Atlas databases.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdatafederationspecstoragedatabasesindex">databases</a></b></td>
        <td>[]object</td>
        <td>
          Array that contains the queryable databases and collections for this data lake.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatafederationspecstoragestoresindex">stores</a></b></td>
        <td>[]object</td>
        <td>
          Array that contains the data stores for the data lake.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec.storage.databases[index]
<sup><sup>[↩ Parent](#atlasdatafederationspecstorage)</sup></sup>



Database associated with this data lake. Databases contain collections and views.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdatafederationspecstoragedatabasesindexcollectionsindex">collections</a></b></td>
        <td>[]object</td>
        <td>
          Array of collections and data sources that map to a stores data store.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxWildcardCollections</b></td>
        <td>integer</td>
        <td>
          Maximum number of wildcard collections in the database. This only applies to S3 data sources.
Minimum value is 1, maximum value is 1000. Default value is 100.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the database to which the data lake maps data.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdatafederationspecstoragedatabasesindexviewsindex">views</a></b></td>
        <td>[]object</td>
        <td>
          Array of aggregation pipelines that apply to the collection. This only applies to S3 data sources.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec.storage.databases[index].collections[index]
<sup><sup>[↩ Parent](#atlasdatafederationspecstoragedatabasesindex)</sup></sup>



Collection maps to a stores data store.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdatafederationspecstoragedatabasesindexcollectionsindexdatasourcesindex">dataSources</a></b></td>
        <td>[]object</td>
        <td>
          Array that contains the data stores that map to a collection for this data lake.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the collection to which MongoDB Atlas maps the data in the data stores.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec.storage.databases[index].collections[index].dataSources[index]
<sup><sup>[↩ Parent](#atlasdatafederationspecstoragedatabasesindexcollectionsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>allowInsecure</b></td>
        <td>boolean</td>
        <td>
          Flag that validates the scheme in the specified URLs.
If true, allows insecure HTTP scheme, doesn't verify the server's certificate chain and hostname, and accepts any certificate with any hostname presented by the server.
If false, allows secure HTTPS scheme only.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>collection</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the collection in the database. For creating a wildcard (*) collection, you must omit this parameter.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>collectionRegex</b></td>
        <td>string</td>
        <td>
          Regex pattern to use for creating the wildcard (*) collection.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the database, which contains the collection in the cluster. You must omit this parameter to generate wildcard (*) collections for dynamically generated databases.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>databaseRegex</b></td>
        <td>string</td>
        <td>
          Regex pattern to use for creating the wildcard (*) database.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>defaultFormat</b></td>
        <td>enum</td>
        <td>
          File format that MongoDB Cloud uses if it encounters a file without a file extension while searching storeName.<br/>
          <br/>
            <i>Enum</i>: .avro, .avro.bz2, .avro.gz, .bson, .bson.bz2, .bson.gz, .bsonx, .csv, .csv.bz2, .csv.gz, .json, .json.bz2, .json.gz, .orc, .parquet, .tsv, .tsv.bz2, .tsv.gz<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          File path that controls how MongoDB Cloud searches for and parses files in the storeName before mapping them to a collection.
Specify / to capture all files and folders from the prefix path.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>provenanceFieldName</b></td>
        <td>string</td>
        <td>
          Name for the field that includes the provenance of the documents in the results. MongoDB Atlas returns different fields in the results for each supported provider.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storeName</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the data store that MongoDB Cloud maps to the collection.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>urls</b></td>
        <td>[]string</td>
        <td>
          URLs of the publicly accessible data files. You can't specify URLs that require authentication.
Atlas Data Lake creates a partition for each URL. If empty or omitted, Data Lake uses the URLs from the store specified in the storeName parameter.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec.storage.databases[index].views[index]
<sup><sup>[↩ Parent](#atlasdatafederationspecstoragedatabasesindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the view, which corresponds to an aggregation pipeline on a collection.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>pipeline</b></td>
        <td>string</td>
        <td>
          Aggregation pipeline stages to apply to the source collection.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>source</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the source collection for the view.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.spec.storage.stores[index]
<sup><sup>[↩ Parent](#atlasdatafederationspecstorage)</sup></sup>



Store is a group of settings that define where the data is stored.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>additionalStorageClasses</b></td>
        <td>[]string</td>
        <td>
          Collection of AWS S3 storage classes. Atlas Data Lake includes the files in these storage classes in the query results.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>bucket</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the AWS S3 bucket.
This label must exactly match the name of an S3 bucket that the data lake can access with the configured AWS Identity and Access Management (IAM) credentials.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>delimiter</b></td>
        <td>string</td>
        <td>
          The delimiter that separates path segments in the data store.
MongoDB Atlas uses the delimiter to efficiently traverse S3 buckets with a hierarchical directory structure. You can specify any character supported by the S3 object keys as the delimiter.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>includeTags</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to use S3 tags on the files in the given path as additional partition attributes.
If set to true, data lake adds the S3 tags as additional partition attributes and adds new top-level BSON elements associating each tag to each document.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the data store. The storeName field references this values as part of the mapping configuration.
To use MongoDB Atlas as a data store, the data lake requires a serverless instance or an M10 or higher cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>prefix</b></td>
        <td>string</td>
        <td>
          Prefix that MongoDB Cloud applies when searching for files in the S3 bucket.
The data store prepends the value of prefix to the path to create the full path for files to ingest.
If omitted, MongoDB Cloud searches all files from the root of the S3 bucket.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>provider</b></td>
        <td>string</td>
        <td>
          The provider used for data stores.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>public</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether the bucket is public.
If set to true, MongoDB Cloud doesn't use the configured AWS Identity and Access Management (IAM) role to access the S3 bucket.
If set to false, the configured AWS IAM role must include permissions to access the S3 bucket.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          Physical location where MongoDB Cloud deploys your AWS-hosted MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases.
When MongoDB Atlas deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Atlas creates them as part of the deployment.
To limit a new VPC peering connection to one CIDR block and region, create the connection first. Deploy the cluster after the connection starts.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.status
<sup><sup>[↩ Parent](#atlasdatafederation)</sup></sup>



DataFederationStatus defines the observed state of AtlasDataFederation.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdatafederationstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mongoDBVersion</b></td>
        <td>string</td>
        <td>
          MongoDBVersion is the version of MongoDB the cluster runs, in <major version>.<minor version> format.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDataFederation.status.conditions[index]
<sup><sup>[↩ Parent](#atlasdatafederationstatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasDeployment
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasDeployment is the Schema for the atlasdeployments API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasDeployment</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasDeploymentSpec defines the desired state of AtlasDeployment.
Only one of DeploymentSpec, AdvancedDeploymentSpec and ServerlessSpec should be defined.<br/>
          <br/>
            <i>Validations</i>:<li>(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef)): must define only one project reference through externalProjectRef or projectRef</li><li>(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef): must define a local connection secret when referencing an external project</li><li>!has(self.serverlessSpec) || (oldSelf.hasValue() && oldSelf.value().serverlessSpec != null): serverlessSpec cannot be added - serverless instances are deprecated</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentstatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasDeploymentStatus defines the observed state of AtlasDeployment.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec
<sup><sup>[↩ Parent](#atlasdeployment)</sup></sup>



AtlasDeploymentSpec defines the desired state of AtlasDeployment.
Only one of DeploymentSpec, AdvancedDeploymentSpec and ServerlessSpec should be defined.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdeploymentspecbackupref">backupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to the backup schedule for the AtlasDeployment.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecconnectionsecret">connectionSecret</a></b></td>
        <td>object</td>
        <td>
          Name of the secret containing Atlas API private and public keys.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspec">deploymentSpec</a></b></td>
        <td>object</td>
        <td>
          Configuration for the advanced (v1.5) deployment API https://www.mongodb.com/docs/atlas/reference/api/clusters/<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecexternalprojectref">externalProjectRef</a></b></td>
        <td>object</td>
        <td>
          externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecflexspec">flexSpec</a></b></td>
        <td>object</td>
        <td>
          Configuration for the Flex cluster API. https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Flex-Clusters<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecprocessargs">processArgs</a></b></td>
        <td>object</td>
        <td>
          ProcessArgs allows modification of Advanced Configuration Options.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecprojectref">projectRef</a></b></td>
        <td>object</td>
        <td>
          projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecserverlessspec">serverlessSpec</a></b></td>
        <td>object</td>
        <td>
          Configuration for the serverless deployment API. https://www.mongodb.com/docs/atlas/reference/api/serverless-instances/
DEPRECATED: Serverless instances are deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>upgradeToDedicated</b></td>
        <td>boolean</td>
        <td>
           upgradeToDedicated, when set to true, triggers the migration from a Flex to a
 Dedicated cluster. The user MUST provide the new dedicated cluster configuration.
 This flag is ignored if the cluster is already dedicated.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.backupRef
<sup><sup>[↩ Parent](#atlasdeploymentspec)</sup></sup>



Reference to the backup schedule for the AtlasDeployment.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.connectionSecret
<sup><sup>[↩ Parent](#atlasdeploymentspec)</sup></sup>



Name of the secret containing Atlas API private and public keys.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec
<sup><sup>[↩ Parent](#atlasdeploymentspec)</sup></sup>



Configuration for the advanced (v1.5) deployment API https://www.mongodb.com/docs/atlas/reference/api/clusters/

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the advanced deployment as it appears in Atlas.
After Atlas creates the deployment, you can't change its name.
Can only contain ASCII letters, numbers, and hyphens.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: Name cannot be modified after deployment creation</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>backupEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates if the deployment uses Cloud Backups for backups.
Applicable only for M10+ deployments.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecbiconnector">biConnector</a></b></td>
        <td>object</td>
        <td>
          Configuration of BI Connector for Atlas on this deployment.
The MongoDB Connector for Business Intelligence for Atlas (BI Connector) is only available for M10 and larger deployments.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>clusterType</b></td>
        <td>enum</td>
        <td>
          Type of the deployment that you want to create.
The parameter is required if replicationSpecs are set or if Global Deployments are deployed.<br/>
          <br/>
            <i>Enum</i>: REPLICASET, SHARDED, GEOSHARDED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configServerManagementMode</b></td>
        <td>enum</td>
        <td>
          Config Server Management Mode for creating or updating a sharded cluster.<br/>
          <br/>
            <i>Enum</i>: ATLAS_MANAGED, FIXED_TO_DEDICATED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspeccustomzonemappingindex">customZoneMapping</a></b></td>
        <td>[]object</td>
        <td>
          List that contains Global Cluster parameters that map zones to geographic regions.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>diskSizeGB</b></td>
        <td>integer</td>
        <td>
          Capacity, in gigabytes, of the host's root volume.
Increase this number to add capacity, up to a maximum possible value of 4096 (i.e., 4 TB).
This value must be a positive integer.
The parameter is required if replicationSpecs are configured.<br/>
          <br/>
            <i>Minimum</i>: 0<br/>
            <i>Maximum</i>: 4096<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>encryptionAtRestProvider</b></td>
        <td>enum</td>
        <td>
          Cloud service provider that offers Encryption at Rest.<br/>
          <br/>
            <i>Enum</i>: AWS, GCP, AZURE, NONE<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspeclabelsindex">labels</a></b></td>
        <td>[]object</td>
        <td>
          Collection of key-value pairs that tag and categorize the deployment.
Each key and value has a maximum length of 255 characters.
DEPRECATED: Cluster labels are deprecated and will be removed in a future release. We strongly recommend that you use Resource Tags instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecmanagednamespacesindex">managedNamespaces</a></b></td>
        <td>[]object</td>
        <td>
          List that contains information to create a managed namespace in a specified Global Cluster to create.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mongoDBMajorVersion</b></td>
        <td>string</td>
        <td>
          MongoDB major version of the cluster. Set to the binary major version.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mongoDBVersion</b></td>
        <td>string</td>
        <td>
          Version of MongoDB that the cluster runs.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>paused</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether the deployment should be paused.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>pitEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates the deployment uses continuous cloud backups.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecreplicationspecsindex">replicationSpecs</a></b></td>
        <td>[]object</td>
        <td>
          Configuration for deployment regions.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>rootCertType</b></td>
        <td>string</td>
        <td>
          Root Certificate Authority that MongoDB Atlas cluster uses.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecsearchindexesindex">searchIndexes</a></b></td>
        <td>[]object</td>
        <td>
          An array of SearchIndex objects with fields that describe the search index.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecsearchnodesindex">searchNodes</a></b></td>
        <td>[]object</td>
        <td>
          Settings for Search Nodes for the cluster. Currently, at most one search node configuration may be defined.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspectagsindex">tags</a></b></td>
        <td>[]object</td>
        <td>
          Key-value pairs for resource tagging.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationProtectionEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether termination protection is enabled on the cluster. If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>versionReleaseSystem</b></td>
        <td>string</td>
        <td>
          Method by which the cluster maintains the MongoDB versions.
If value is CONTINUOUS, you must not specify mongoDBMajorVersion.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.biConnector
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspec)</sup></sup>



Configuration of BI Connector for Atlas on this deployment.
The MongoDB Connector for Business Intelligence for Atlas (BI Connector) is only available for M10 and larger deployments.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether the Business Intelligence Connector for Atlas is enabled on the deployment.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readPreference</b></td>
        <td>string</td>
        <td>
          Source from which the BI Connector for Atlas reads data. Each BI Connector for Atlas read preference contains a distinct combination of readPreference and readPreferenceTags options.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.customZoneMapping[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>location</b></td>
        <td>string</td>
        <td>
          Code that represents a location that maps to a zone in your global cluster.
MongoDB Atlas represents this location with a ISO 3166-2 location and subdivision codes when possible.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>zone</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the zone in your global cluster. This zone maps to a location code.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.labels[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspec)</sup></sup>



LabelSpec contains key-value pairs that tag and categorize the Cluster/DBUser

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Key applied to tag and categorize this component.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Value set to the Key applied to tag and categorize this component.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.managedNamespaces[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspec)</sup></sup>



ManagedNamespace represents the information about managed namespace configuration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>collection</b></td>
        <td>string</td>
        <td>
          Human-readable label of the collection to manage for this Global Cluster.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>db</b></td>
        <td>string</td>
        <td>
          Human-readable label of the database to manage for this Global Cluster.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>customShardKey</b></td>
        <td>string</td>
        <td>
          Database parameter used to divide the collection into shards. Global clusters require a compound shard key.
This compound shard key combines the location parameter and the user-selected custom key.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isCustomShardKeyHashed</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether someone hashed the custom shard key for the specified collection.
If you set this value to false, MongoDB Cloud uses ranged sharding.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isShardKeyUnique</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether someone hashed the custom shard key.
If this parameter returns false, this cluster uses ranged sharding.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>numInitialChunks</b></td>
        <td>integer</td>
        <td>
          Minimum number of chunks to create initially when sharding an empty collection with a hashed shard key.
Maximum value is 8192.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>presplitHashedZones</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether MongoDB Cloud should create and distribute initial chunks for an empty or non-existing collection.
MongoDB Cloud distributes data based on the defined zones and zone ranges for the collection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.replicationSpecs[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>numShards</b></td>
        <td>integer</td>
        <td>
          Positive integer that specifies the number of shards to deploy in each specified zone.
If you set this value to 1 and clusterType is SHARDED, MongoDB Cloud deploys a single-shard sharded cluster.
Don't create a sharded cluster with a single shard for production environments.
Single-shard sharded clusters don't provide the same benefits as multi-shard configurations<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindex">regionConfigs</a></b></td>
        <td>[]object</td>
        <td>
          Hardware specifications for nodes set for a given region.
Each regionConfigs object describes the region's priority in elections and the number and type of MongoDB nodes that MongoDB Cloud deploys to the region.
Each regionConfigs object must have either an analyticsSpecs object, electableSpecs object, or readOnlySpecs object.
Tenant clusters only require electableSpecs. Dedicated clusters can specify any of these specifications, but must have at least one electableSpecs object within a replicationSpec.
Every hardware specification must use the same instanceSize.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>zoneName</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the zone in a Global Cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.replicationSpecs[index].regionConfigs[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecreplicationspecsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindexanalyticsspecs">analyticsSpecs</a></b></td>
        <td>object</td>
        <td>
          Hardware specifications for analytics nodes deployed in the region.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindexautoscaling">autoScaling</a></b></td>
        <td>object</td>
        <td>
          Options that determine how this cluster handles resource scaling.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>backingProviderName</b></td>
        <td>enum</td>
        <td>
          Cloud service provider on which the host for a multi-tenant deployment is provisioned.
This setting only works when "providerName" : "TENANT" and "providerSetting.instanceSizeName" : M2 or M5.
Otherwise, it should be equal to the "providerName" value.<br/>
          <br/>
            <i>Enum</i>: AWS, GCP, AZURE<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindexelectablespecs">electableSpecs</a></b></td>
        <td>object</td>
        <td>
          Hardware specifications for nodes deployed in the region.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>priority</b></td>
        <td>integer</td>
        <td>
          Precedence is given to this region when a primary election occurs.
If your regionConfigs has only readOnlySpecs, analyticsSpecs, or both, set this value to 0.
If you have multiple regionConfigs objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order.
The highest priority is 7<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>providerName</b></td>
        <td>enum</td>
        <td>
          <br/>
          <br/>
            <i>Enum</i>: AWS, GCP, AZURE, TENANT, SERVERLESS<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindexreadonlyspecs">readOnlySpecs</a></b></td>
        <td>object</td>
        <td>
          Hardware specifications for read only nodes deployed in the region.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>regionName</b></td>
        <td>string</td>
        <td>
          Physical location of your MongoDB deployment.
The region you choose can affect network latency for clients accessing your databases.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.replicationSpecs[index].regionConfigs[index].analyticsSpecs
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindex)</sup></sup>



Hardware specifications for analytics nodes deployed in the region.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>diskIOPS</b></td>
        <td>integer</td>
        <td>
          Disk IOPS setting for AWS storage.
Set only if you selected AWS as your cloud service provider.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ebsVolumeType</b></td>
        <td>enum</td>
        <td>
          Disk IOPS setting for AWS storage.
Set only if you selected AWS as your cloud service provider.<br/>
          <br/>
            <i>Enum</i>: STANDARD, PROVISIONED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>instanceSize</b></td>
        <td>string</td>
        <td>
          Hardware specification for the instance sizes in this region.
Each instance size has a default storage and memory capacity.
The instance size you select applies to all the data-bearing hosts in your instance size.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>nodeCount</b></td>
        <td>integer</td>
        <td>
          Number of nodes of the given type for MongoDB Cloud to deploy to the region.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.replicationSpecs[index].regionConfigs[index].autoScaling
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindex)</sup></sup>



Options that determine how this cluster handles resource scaling.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindexautoscalingcompute">compute</a></b></td>
        <td>object</td>
        <td>
          Collection of settings that configure how a deployment might scale its deployment tier and whether the deployment can scale down.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindexautoscalingdiskgb">diskGB</a></b></td>
        <td>object</td>
        <td>
          Flag that indicates whether disk auto-scaling is enabled. The default is true.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.replicationSpecs[index].regionConfigs[index].autoScaling.compute
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindexautoscaling)</sup></sup>



Collection of settings that configure how a deployment might scale its deployment tier and whether the deployment can scale down.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether deployment tier auto-scaling is enabled. The default is false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxInstanceSize</b></td>
        <td>string</td>
        <td>
          Maximum instance size to which your deployment can automatically scale (such as M40). Atlas requires this parameter if "autoScaling.compute.enabled" : true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>minInstanceSize</b></td>
        <td>string</td>
        <td>
          Minimum instance size to which your deployment can automatically scale (such as M10). Atlas requires this parameter if "autoScaling.compute.scaleDownEnabled" : true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scaleDownEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether the deployment tier may scale down. Atlas requires this parameter if "autoScaling.compute.enabled" : true.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.replicationSpecs[index].regionConfigs[index].autoScaling.diskGB
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindexautoscaling)</sup></sup>



Flag that indicates whether disk auto-scaling is enabled. The default is true.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether this cluster enables disk auto-scaling.
The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.replicationSpecs[index].regionConfigs[index].electableSpecs
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindex)</sup></sup>



Hardware specifications for nodes deployed in the region.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>diskIOPS</b></td>
        <td>integer</td>
        <td>
          Disk IOPS setting for AWS storage.
Set only if you selected AWS as your cloud service provider.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ebsVolumeType</b></td>
        <td>enum</td>
        <td>
          Disk IOPS setting for AWS storage.
Set only if you selected AWS as your cloud service provider.<br/>
          <br/>
            <i>Enum</i>: STANDARD, PROVISIONED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>instanceSize</b></td>
        <td>string</td>
        <td>
          Hardware specification for the instance sizes in this region.
Each instance size has a default storage and memory capacity.
The instance size you select applies to all the data-bearing hosts in your instance size.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>nodeCount</b></td>
        <td>integer</td>
        <td>
          Number of nodes of the given type for MongoDB Cloud to deploy to the region.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.replicationSpecs[index].regionConfigs[index].readOnlySpecs
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecreplicationspecsindexregionconfigsindex)</sup></sup>



Hardware specifications for read only nodes deployed in the region.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>diskIOPS</b></td>
        <td>integer</td>
        <td>
          Disk IOPS setting for AWS storage.
Set only if you selected AWS as your cloud service provider.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ebsVolumeType</b></td>
        <td>enum</td>
        <td>
          Disk IOPS setting for AWS storage.
Set only if you selected AWS as your cloud service provider.<br/>
          <br/>
            <i>Enum</i>: STANDARD, PROVISIONED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>instanceSize</b></td>
        <td>string</td>
        <td>
          Hardware specification for the instance sizes in this region.
Each instance size has a default storage and memory capacity.
The instance size you select applies to all the data-bearing hosts in your instance size.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>nodeCount</b></td>
        <td>integer</td>
        <td>
          Number of nodes of the given type for MongoDB Cloud to deploy to the region.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.searchIndexes[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspec)</sup></sup>



SearchIndex is the CRD to configure part of the Atlas Search Index.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>DBName</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the database that contains the collection with one or more Atlas Search indexes.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>collectionName</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the collection that contains one or more Atlas Search indexes.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies this index. Must be unique for a deployment.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Type of the index.<br/>
          <br/>
            <i>Enum</i>: search, vectorSearch<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecsearchindexesindexsearch">search</a></b></td>
        <td>object</td>
        <td>
          Atlas search index configuration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecsearchindexesindexvectorsearch">vectorSearch</a></b></td>
        <td>object</td>
        <td>
          Atlas vector search index configuration.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.searchIndexes[index].search
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecsearchindexesindex)</sup></sup>



Atlas search index configuration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecsearchindexesindexsearchmappings">mappings</a></b></td>
        <td>object</td>
        <td>
          Index specifications for the collection's fields.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecsearchindexesindexsearchsearchconfigurationref">searchConfigurationRef</a></b></td>
        <td>object</td>
        <td>
          A reference to the AtlasSearchIndexConfig custom resource.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecsearchindexesindexsearchsynonymsindex">synonyms</a></b></td>
        <td>[]object</td>
        <td>
          Rule sets that map words to their synonyms in this index.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.searchIndexes[index].search.mappings
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecsearchindexesindexsearch)</sup></sup>



Index specifications for the collection's fields.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>dynamic</b></td>
        <td>JSON</td>
        <td>
          Indicates whether the index uses static, default dynamic, or configurable dynamic mappings.
Set to **true** to enable dynamic mapping with default type set or define object to specify the name of the configured type sets for dynamic mapping.
If you specify configurable dynamic mappings, you must define the referred type sets in the **typeSets** field.
Set to **false** to use only static mappings through **mappings.fields**.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>fields</b></td>
        <td>JSON</td>
        <td>
          One or more field specifications for the Atlas Search index. Required if mapping.dynamic is omitted or set to false.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.searchIndexes[index].search.searchConfigurationRef
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecsearchindexesindexsearch)</sup></sup>



A reference to the AtlasSearchIndexConfig custom resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.searchIndexes[index].search.synonyms[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecsearchindexesindexsearch)</sup></sup>



Synonym represents "Synonym" type of Atlas Search Index.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>analyzer</b></td>
        <td>enum</td>
        <td>
          Specific pre-defined method chosen to apply to the synonyms to be searched.<br/>
          <br/>
            <i>Enum</i>: lucene.standard, lucene.simple, lucene.whitespace, lucene.keyword, lucene.arabic, lucene.armenian, lucene.basque, lucene.bengali, lucene.brazilian, lucene.bulgarian, lucene.catalan, lucene.chinese, lucene.cjk, lucene.czech, lucene.danish, lucene.dutch, lucene.english, lucene.finnish, lucene.french, lucene.galician, lucene.german, lucene.greek, lucene.hindi, lucene.hungarian, lucene.indonesian, lucene.irish, lucene.italian, lucene.japanese, lucene.korean, lucene.kuromoji, lucene.latvian, lucene.lithuanian, lucene.morfologik, lucene.nori, lucene.norwegian, lucene.persian, lucene.portuguese, lucene.romanian, lucene.russian, lucene.smartcn, lucene.sorani, lucene.spanish, lucene.swedish, lucene.thai, lucene.turkish, lucene.ukrainian<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the synonym definition. Each name must be unique within the same index definition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecdeploymentspecsearchindexesindexsearchsynonymsindexsource">source</a></b></td>
        <td>object</td>
        <td>
          Data set that stores the mapping one or more words map to one or more synonyms of those words.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.searchIndexes[index].search.synonyms[index].source
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecsearchindexesindexsearchsynonymsindex)</sup></sup>



Data set that stores the mapping one or more words map to one or more synonyms of those words.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>collection</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the MongoDB collection that stores words and their applicable synonyms.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.searchIndexes[index].vectorSearch
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspecsearchindexesindex)</sup></sup>



Atlas vector search index configuration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fields</b></td>
        <td>JSON</td>
        <td>
          Array of JSON objects. See examples https://dochub.mongodb.org/core/avs-vector-type<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.searchNodes[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>instanceSize</b></td>
        <td>enum</td>
        <td>
          Hardware specification for the Search Node instance sizes.<br/>
          <br/>
            <i>Enum</i>: S20_HIGHCPU_NVME, S30_HIGHCPU_NVME, S40_HIGHCPU_NVME, S50_HIGHCPU_NVME, S60_HIGHCPU_NVME, S70_HIGHCPU_NVME, S80_HIGHCPU_NVME, S30_LOWCPU_NVME, S40_LOWCPU_NVME, S50_LOWCPU_NVME, S60_LOWCPU_NVME, S80_LOWCPU_NVME, S90_LOWCPU_NVME, S100_LOWCPU_NVME, S110_LOWCPU_NVME<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>nodeCount</b></td>
        <td>integer</td>
        <td>
          Number of Search Nodes in the cluster.<br/>
          <br/>
            <i>Minimum</i>: 2<br/>
            <i>Maximum</i>: 32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.deploymentSpec.tags[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecdeploymentspec)</sup></sup>



TagSpec holds a key-value pair for resource tagging on this deployment.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Constant that defines the set of the tag.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Variable that belongs to the set of the tag.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.externalProjectRef
<sup><sup>[↩ Parent](#atlasdeploymentspec)</sup></sup>



externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID is the Atlas project ID.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.flexSpec
<sup><sup>[↩ Parent](#atlasdeploymentspec)</sup></sup>



Configuration for the Flex cluster API. https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Flex-Clusters

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the instance.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecflexspecprovidersettings">providerSettings</a></b></td>
        <td>object</td>
        <td>
          Group of cloud provider settings that configure the provisioned MongoDB flex cluster.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecflexspectagsindex">tags</a></b></td>
        <td>[]object</td>
        <td>
          List that contains key-value pairs between 1 and 255 characters in length for tagging and categorizing the instance.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationProtectionEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether termination protection is enabled on the cluster.
If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.flexSpec.providerSettings
<sup><sup>[↩ Parent](#atlasdeploymentspecflexspec)</sup></sup>



Group of cloud provider settings that configure the provisioned MongoDB flex cluster.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>backingProviderName</b></td>
        <td>enum</td>
        <td>
          Cloud service provider on which MongoDB Atlas provisions the flex cluster.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: Backing Provider cannot be modified after cluster creation</li>
            <i>Enum</i>: AWS, GCP, AZURE<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>regionName</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the geographic location of your MongoDB flex cluster.
The region you choose can affect network latency for clients accessing your databases.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: Region Name cannot be modified after cluster creation</li>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.flexSpec.tags[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecflexspec)</sup></sup>



TagSpec holds a key-value pair for resource tagging on this deployment.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Constant that defines the set of the tag.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Variable that belongs to the set of the tag.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.processArgs
<sup><sup>[↩ Parent](#atlasdeploymentspec)</sup></sup>



ProcessArgs allows modification of Advanced Configuration Options.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>defaultReadConcern</b></td>
        <td>string</td>
        <td>
          String that indicates the default level of acknowledgment requested from MongoDB for read operations set for this cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>defaultWriteConcern</b></td>
        <td>string</td>
        <td>
          String that indicates the default level of acknowledgment requested from MongoDB for write operations set for this cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failIndexKeyTooLong</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to fail the operation and return an error when you insert or update documents where all indexed entries exceed 1024 bytes.
If you set this to false, mongod writes documents that exceed this limit, but doesn't index them.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>javascriptEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether the cluster allows execution of operations that perform server-side executions of JavaScript.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>minimumEnabledTlsProtocol</b></td>
        <td>string</td>
        <td>
          String that indicates the minimum TLS version that the cluster accepts for incoming connections.
Clusters using TLS 1.0 or 1.1 should consider setting TLS 1.2 as the minimum TLS protocol version.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>noTableScan</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether the cluster disables executing any query that requires a collection scan to return results.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>oplogMinRetentionHours</b></td>
        <td>string</td>
        <td>
          Minimum retention window for cluster's oplog expressed in hours. A value of null indicates that the cluster uses the default minimum oplog window that MongoDB Cloud calculates.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>oplogSizeMB</b></td>
        <td>integer</td>
        <td>
          Number that indicates the storage limit of a cluster's oplog expressed in megabytes.
A value of null indicates that the cluster uses the default oplog size that Atlas calculates.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sampleRefreshIntervalBIConnector</b></td>
        <td>integer</td>
        <td>
          Number that indicates the documents per database to sample when gathering schema information.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sampleSizeBIConnector</b></td>
        <td>integer</td>
        <td>
          Number that indicates the interval in seconds at which the mongosqld process re-samples data to create its relational schema.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.projectRef
<sup><sup>[↩ Parent](#atlasdeploymentspec)</sup></sup>



projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.serverlessSpec
<sup><sup>[↩ Parent](#atlasdeploymentspec)</sup></sup>



Configuration for the serverless deployment API. https://www.mongodb.com/docs/atlas/reference/api/serverless-instances/
DEPRECATED: Serverless instances are deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the serverless deployment as it appears in Atlas.
After Atlas creates the deployment, you can't change its name.
Can only contain ASCII letters, numbers, and hyphens.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecserverlessspecprovidersettings">providerSettings</a></b></td>
        <td>object</td>
        <td>
          Configuration for the provisioned hosts on which MongoDB runs. The available options are specific to the cloud service provider.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecserverlessspecbackupoptions">backupOptions</a></b></td>
        <td>object</td>
        <td>
          Serverless Backup Options<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecserverlessspecprivateendpointsindex">privateEndpoints</a></b></td>
        <td>[]object</td>
        <td>
          List that contains the private endpoint configurations for the Serverless instance.
DEPRECATED: Serverless private endpoints are deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecserverlessspectagsindex">tags</a></b></td>
        <td>[]object</td>
        <td>
          Key-value pairs for resource tagging.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationProtectionEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether termination protection is enabled on the cluster. If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.serverlessSpec.providerSettings
<sup><sup>[↩ Parent](#atlasdeploymentspecserverlessspec)</sup></sup>



Configuration for the provisioned hosts on which MongoDB runs. The available options are specific to the cloud service provider.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>providerName</b></td>
        <td>enum</td>
        <td>
          Cloud service provider on which Atlas provisions the hosts.<br/>
          <br/>
            <i>Enum</i>: AWS, GCP, AZURE, TENANT, SERVERLESS<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecserverlessspecprovidersettingsautoscaling">autoScaling</a></b></td>
        <td>object</td>
        <td>
          Range of instance sizes to which your deployment can scale.
DEPRECATED: The value of this field doesn't take any effect.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>backingProviderName</b></td>
        <td>enum</td>
        <td>
          Cloud service provider on which the host for a multi-tenant deployment is provisioned.
This setting only works when "providerSetting.providerName" : "TENANT" and "providerSetting.instanceSizeName" : M2 or M5.<br/>
          <br/>
            <i>Enum</i>: AWS, GCP, AZURE<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>diskIOPS</b></td>
        <td>integer</td>
        <td>
          Disk IOPS setting for AWS storage.
Set only if you selected AWS as your cloud service provider.
DEPRECATED: The value of this field doesn't take any effect.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>diskTypeName</b></td>
        <td>string</td>
        <td>
          Type of disk if you selected Azure as your cloud service provider.
DEPRECATED: The value of this field doesn't take any effect.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>encryptEBSVolume</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether the Amazon EBS encryption feature encrypts the host's root volume for both data at rest within the volume and for data moving between the volume and the deployment.
DEPRECATED: The value of this field doesn't take any effect.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>instanceSizeName</b></td>
        <td>string</td>
        <td>
          Atlas provides different deployment tiers, each with a default storage capacity and RAM size. The deployment you select is used for all the data-bearing hosts in your deployment tier.
DEPRECATED: The value of this field doesn't take any effect.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>regionName</b></td>
        <td>string</td>
        <td>
          Physical location of your MongoDB deployment.
The region you choose can affect network latency for clients accessing your databases.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeType</b></td>
        <td>enum</td>
        <td>
          Disk IOPS setting for AWS storage.
Set only if you selected AWS as your cloud service provider.
DEPRECATED: The value of this field doesn't take any effect.<br/>
          <br/>
            <i>Enum</i>: STANDARD, PROVISIONED<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.serverlessSpec.providerSettings.autoScaling
<sup><sup>[↩ Parent](#atlasdeploymentspecserverlessspecprovidersettings)</sup></sup>



Range of instance sizes to which your deployment can scale.
DEPRECATED: The value of this field doesn't take any effect.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>autoIndexingEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether autopilot mode for Performance Advisor is enabled.
The default is false.
DEPRECATED: This flag is no longer supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentspecserverlessspecprovidersettingsautoscalingcompute">compute</a></b></td>
        <td>object</td>
        <td>
          Collection of settings that configure how a deployment might scale its deployment tier and whether the deployment can scale down.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>diskGBEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether disk auto-scaling is enabled. The default is true.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.serverlessSpec.providerSettings.autoScaling.compute
<sup><sup>[↩ Parent](#atlasdeploymentspecserverlessspecprovidersettingsautoscaling)</sup></sup>



Collection of settings that configure how a deployment might scale its deployment tier and whether the deployment can scale down.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether deployment tier auto-scaling is enabled. The default is false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxInstanceSize</b></td>
        <td>string</td>
        <td>
          Maximum instance size to which your deployment can automatically scale (such as M40). Atlas requires this parameter if "autoScaling.compute.enabled" : true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>minInstanceSize</b></td>
        <td>string</td>
        <td>
          Minimum instance size to which your deployment can automatically scale (such as M10). Atlas requires this parameter if "autoScaling.compute.scaleDownEnabled" : true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scaleDownEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether the deployment tier may scale down. Atlas requires this parameter if "autoScaling.compute.enabled" : true.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.serverlessSpec.backupOptions
<sup><sup>[↩ Parent](#atlasdeploymentspecserverlessspec)</sup></sup>



Serverless Backup Options

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>serverlessContinuousBackupEnabled</b></td>
        <td>boolean</td>
        <td>
          ServerlessContinuousBackupEnabled indicates whether the cluster uses continuous cloud backups.
DEPRECATED: Serverless instances are deprecated, and no longer support continuous backup. See https://dochub.mongodb.org/core/atlas-flex-migration for details.<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.serverlessSpec.privateEndpoints[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecserverlessspec)</sup></sup>



ServerlessPrivateEndpoint configures private endpoints for the Serverless instances.
DEPRECATED: Serverless private endpoints are deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>cloudProviderEndpointID</b></td>
        <td>string</td>
        <td>
          CloudProviderEndpointID is the identifier of the cloud provider endpoint.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of the Serverless PrivateLink Service. Should be unique.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>privateEndpointIpAddress</b></td>
        <td>string</td>
        <td>
          PrivateEndpointIPAddress is the IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.spec.serverlessSpec.tags[index]
<sup><sup>[↩ Parent](#atlasdeploymentspecserverlessspec)</sup></sup>



TagSpec holds a key-value pair for resource tagging on this deployment.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Constant that defines the set of the tag.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Variable that belongs to the set of the tag.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.status
<sup><sup>[↩ Parent](#atlasdeployment)</sup></sup>



AtlasDeploymentStatus defines the observed state of AtlasDeployment.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasdeploymentstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentstatusconnectionstrings">connectionStrings</a></b></td>
        <td>object</td>
        <td>
          ConnectionStrings is a set of connection strings that your applications use to connect to this cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentstatuscustomzonemapping">customZoneMapping</a></b></td>
        <td>object</td>
        <td>
          List that contains key value pairs to map zones to geographic regions.
These pairs map an ISO 3166-1a2 location code, with an ISO 3166-2 subdivision code when possible, to a unique 24-hexadecimal string that identifies the custom zone.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentstatusmanagednamespacesindex">managedNamespaces</a></b></td>
        <td>[]object</td>
        <td>
          List that contains a namespace for a Global Cluster. MongoDB Atlas manages this cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mongoDBVersion</b></td>
        <td>string</td>
        <td>
          MongoDBVersion is the version of MongoDB the cluster runs, in <major version>.<minor version> format.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mongoURIUpdated</b></td>
        <td>string</td>
        <td>
          MongoURIUpdated is a timestamp in ISO 8601 date and time format in UTC when the connection string was last updated.
The connection string changes if you update any of the other values.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentstatusreplicasetsindex">replicaSets</a></b></td>
        <td>[]object</td>
        <td>
          Details that explain how MongoDB Cloud replicates data on the specified MongoDB database.
This array has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentstatussearchindexesindex">searchIndexes</a></b></td>
        <td>[]object</td>
        <td>
          SearchIndexes contains a list of search indexes statuses configured for a project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentstatusserverlessprivateendpointsindex">serverlessPrivateEndpoints</a></b></td>
        <td>[]object</td>
        <td>
          ServerlessPrivateEndpoints contains a list of private endpoints configured for the serverless deployment.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stateName</b></td>
        <td>string</td>
        <td>
          StateName is the current state of the cluster.
The possible states are: IDLE, CREATING, UPDATING, DELETING, DELETED, REPAIRING<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.status.conditions[index]
<sup><sup>[↩ Parent](#atlasdeploymentstatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.status.connectionStrings
<sup><sup>[↩ Parent](#atlasdeploymentstatus)</sup></sup>



ConnectionStrings is a set of connection strings that your applications use to connect to this cluster.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>private</b></td>
        <td>string</td>
        <td>
          Network-peering-endpoint-aware mongodb:// connection strings for each interface VPC endpoint you configured to connect to this cluster.
Atlas returns this parameter only if you created a network peering connection to this cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentstatusconnectionstringsprivateendpointindex">privateEndpoint</a></b></td>
        <td>[]object</td>
        <td>
          Private endpoint connection strings.
Each object describes the connection strings you can use to connect to this cluster through a private endpoint.
Atlas returns this parameter only if you deployed a private endpoint to all regions to which you deployed this cluster's nodes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>privateSrv</b></td>
        <td>string</td>
        <td>
          Network-peering-endpoint-aware mongodb+srv:// connection strings for each interface VPC endpoint you configured to connect to this cluster.
Atlas returns this parameter only if you created a network peering connection to this cluster.
Use this URI format if your driver supports it. If it doesn't, use connectionStrings.private.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>standard</b></td>
        <td>string</td>
        <td>
          Public mongodb:// connection string for this cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>standardSrv</b></td>
        <td>string</td>
        <td>
          Public mongodb+srv:// connection string for this cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.status.connectionStrings.privateEndpoint[index]
<sup><sup>[↩ Parent](#atlasdeploymentstatusconnectionstrings)</sup></sup>



PrivateEndpoint connection strings. Each object describes the connection strings
you can use to connect to this cluster through a private endpoint.
Atlas returns this parameter only if you deployed a private endpoint to all regions
to which you deployed this cluster's nodes.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>connectionString</b></td>
        <td>string</td>
        <td>
          Private-endpoint-aware mongodb:// connection string for this private endpoint.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasdeploymentstatusconnectionstringsprivateendpointindexendpointsindex">endpoints</a></b></td>
        <td>[]object</td>
        <td>
          Private endpoint through which you connect to Atlas when you use connectionStrings.privateEndpoint[n].connectionString or connectionStrings.privateEndpoint[n].srvConnectionString.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>srvConnectionString</b></td>
        <td>string</td>
        <td>
          Private-endpoint-aware mongodb+srv:// connection string for this private endpoint.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>srvShardOptimizedConnectionString</b></td>
        <td>string</td>
        <td>
          Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of MongoDB process that you connect to with the connection strings

Atlas returns:

• MONGOD for replica sets, or

• MONGOS for sharded clusters<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.status.connectionStrings.privateEndpoint[index].endpoints[index]
<sup><sup>[↩ Parent](#atlasdeploymentstatusconnectionstringsprivateendpointindex)</sup></sup>



Endpoint through which you connect to Atlas

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>endpointId</b></td>
        <td>string</td>
        <td>
          Unique identifier of the private endpoint.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ip</b></td>
        <td>string</td>
        <td>
          Private IP address of the private endpoint network interface you created in your Azure VNet.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>providerName</b></td>
        <td>string</td>
        <td>
          Cloud provider to which you deployed the private endpoint. Atlas returns AWS or AZURE.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          Region to which you deployed the private endpoint.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.status.customZoneMapping
<sup><sup>[↩ Parent](#atlasdeploymentstatus)</sup></sup>



List that contains key value pairs to map zones to geographic regions.
These pairs map an ISO 3166-1a2 location code, with an ISO 3166-2 subdivision code when possible, to a unique 24-hexadecimal string that identifies the custom zone.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>customZoneMapping</b></td>
        <td>map[string]string</td>
        <td>
          List that contains key value pairs to map zones to geographic regions.
These pairs map an ISO 3166-1a2 location code, with an ISO 3166-2 subdivision code when possible, to a unique 24-hexadecimal string that identifies the custom zone.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>zoneMappingErrMessage</b></td>
        <td>string</td>
        <td>
          Error message for failed Custom Zone Mapping.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>zoneMappingState</b></td>
        <td>string</td>
        <td>
          Status of the Custom Zone Mapping.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.status.managedNamespaces[index]
<sup><sup>[↩ Parent](#atlasdeploymentstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>collection</b></td>
        <td>string</td>
        <td>
          Human-readable label of the collection to manage for this Global Cluster.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>db</b></td>
        <td>string</td>
        <td>
          Human-readable label of the database to manage for this Global Cluster.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>customShardKey</b></td>
        <td>string</td>
        <td>
          Database parameter used to divide the collection into shards. Global clusters require a compound shard key.
This compound shard key combines the location parameter and the user-selected custom key.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>errMessage</b></td>
        <td>string</td>
        <td>
          Error message for a failed Managed Namespace.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isCustomShardKeyHashed</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether someone hashed the custom shard key for the specified collection.
If you set this value to false, MongoDB Atlas uses ranged sharding.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isShardKeyUnique</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether someone hashed the custom shard key. If this parameter returns false, this cluster uses ranged sharding.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>numInitialChunks</b></td>
        <td>integer</td>
        <td>
          Minimum number of chunks to create initially when sharding an empty collection with a hashed shard key.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>presplitHashedZones</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether MongoDB Cloud should create and distribute initial chunks for an empty or non-existing collection.
MongoDB Atlas distributes data based on the defined zones and zone ranges for the collection.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the Managed Namespace.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.status.replicaSets[index]
<sup><sup>[↩ Parent](#atlasdeploymentstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Unique 24-hexadecimal digit string that identifies the replication object for a shard in a Cluster.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>zoneName</b></td>
        <td>string</td>
        <td>
          Human-readable label that describes the zone this shard belongs to in a Global Cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasDeployment.status.searchIndexes[index]
<sup><sup>[↩ Parent](#atlasdeploymentstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>ID</b></td>
        <td>string</td>
        <td>
          Unique 24-hexadecimal digit string that identifies this Atlas Search index.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          Details on the status of the search index.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies this index.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Condition of the search index.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasDeployment.status.serverlessPrivateEndpoints[index]
<sup><sup>[↩ Parent](#atlasdeploymentstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>_id</b></td>
        <td>string</td>
        <td>
          ID is the identifier of the Serverless PrivateLink Service.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>cloudProviderEndpointId</b></td>
        <td>string</td>
        <td>
          CloudProviderEndpointID is the identifier of the cloud provider endpoint.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>endpointServiceName</b></td>
        <td>string</td>
        <td>
          EndpointServiceName is the name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>errorMessage</b></td>
        <td>string</td>
        <td>
          ErrorMessage is the error message if the Serverless PrivateLink Service failed to create or connect.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of the Serverless PrivateLink Service. Should be unique.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>privateEndpointIpAddress</b></td>
        <td>string</td>
        <td>
          PrivateEndpointIPAddress is the IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>privateLinkServiceResourceId</b></td>
        <td>string</td>
        <td>
          PrivateLinkServiceResourceID is the root-relative path that identifies the Azure Private Link Service that MongoDB Cloud manages. MongoDB Cloud returns null while it creates the endpoint service.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>providerName</b></td>
        <td>string</td>
        <td>
          ProviderName is human-readable label that identifies the cloud provider. Values include AWS or AZURE.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the AWS Serverless PrivateLink connection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasFederatedAuth
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasFederatedAuth is the Schema for the Atlasfederatedauth API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasFederatedAuth</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasfederatedauthspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasFederatedAuthSpec defines the desired state of AtlasFederatedAuth.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasfederatedauthstatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasFederatedAuthStatus defines the observed state of AtlasFederatedAuth.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasFederatedAuth.spec
<sup><sup>[↩ Parent](#atlasfederatedauth)</sup></sup>



AtlasFederatedAuthSpec defines the desired state of AtlasFederatedAuth.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasfederatedauthspecconnectionsecretref">connectionSecretRef</a></b></td>
        <td>object</td>
        <td>
          Connection secret with API credentials for configuring the federation.
These credentials must have OrganizationOwner permissions.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>dataAccessIdentityProviders</b></td>
        <td>[]string</td>
        <td>
          The collection of unique ids representing the identity providers that can be used for data access in this organization.
Currently connected data access identity providers missing from this field will be disconnected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>domainAllowList</b></td>
        <td>[]string</td>
        <td>
          Approved domains that restrict users who can join the organization based on their email address.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>domainRestrictionEnabled</b></td>
        <td>boolean</td>
        <td>
          Prevent users in the federation from accessing organizations outside the federation, and creating new organizations.
This option applies to the entire federation.
See more information at https://www.mongodb.com/docs/atlas/security/federation-advanced-options/#restrict-user-membership-to-the-federation<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>postAuthRoleGrants</b></td>
        <td>[]string</td>
        <td>
          Atlas roles that are granted to a user in this organization after authenticating.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasfederatedauthspecrolemappingsindex">roleMappings</a></b></td>
        <td>[]object</td>
        <td>
          Map IDP groups to Atlas roles.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ssoDebugEnabled</b></td>
        <td>boolean</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasFederatedAuth.spec.connectionSecretRef
<sup><sup>[↩ Parent](#atlasfederatedauthspec)</sup></sup>



Connection secret with API credentials for configuring the federation.
These credentials must have OrganizationOwner permissions.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasFederatedAuth.spec.roleMappings[index]
<sup><sup>[↩ Parent](#atlasfederatedauthspec)</sup></sup>



RoleMapping maps an external group from an identity provider to roles within Atlas.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>externalGroupName</b></td>
        <td>string</td>
        <td>
          ExternalGroupName is the name of the IDP group to which this mapping applies.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasfederatedauthspecrolemappingsindexroleassignmentsindex">roleAssignments</a></b></td>
        <td>[]object</td>
        <td>
          RoleAssignments define the roles within projects that should be given to members of the group.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasFederatedAuth.spec.roleMappings[index].roleAssignments[index]
<sup><sup>[↩ Parent](#atlasfederatedauthspecrolemappingsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>projectName</b></td>
        <td>string</td>
        <td>
          The Atlas project in the same org in which the role should be given.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>role</b></td>
        <td>enum</td>
        <td>
          The role in Atlas that should be given to group members.<br/>
          <br/>
            <i>Enum</i>: ORG_MEMBER, ORG_READ_ONLY, ORG_BILLING_ADMIN, ORG_GROUP_CREATOR, ORG_OWNER, ORG_BILLING_READ_ONLY, GROUP_OWNER, GROUP_READ_ONLY, GROUP_DATA_ACCESS_ADMIN, GROUP_DATA_ACCESS_READ_ONLY, GROUP_DATA_ACCESS_READ_WRITE, GROUP_CLUSTER_MANAGER, GROUP_SEARCH_INDEX_EDITOR, GROUP_DATABASE_ACCESS_ADMIN, GROUP_BACKUP_MANAGER, GROUP_STREAM_PROCESSING_OWNER, ORG_STREAM_PROCESSING_ADMIN, GROUP_OBSERVABILITY_VIEWER<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasFederatedAuth.status
<sup><sup>[↩ Parent](#atlasfederatedauth)</sup></sup>



AtlasFederatedAuthStatus defines the observed state of AtlasFederatedAuth.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasfederatedauthstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasFederatedAuth.status.conditions[index]
<sup><sup>[↩ Parent](#atlasfederatedauthstatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasIPAccessList
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasIPAccessList is the Schema for the atlasipaccesslists API.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasIPAccessList</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasipaccesslistspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasIPAccessListSpec defines the desired state of AtlasIPAccessList.<br/>
          <br/>
            <i>Validations</i>:<li>(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef)): must define only one project reference through externalProjectRef or projectRef</li><li>(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef): must define a local connection secret when referencing an external project</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasipaccessliststatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasIPAccessListStatus is the most recent observed status of the AtlasIPAccessList cluster. Read-only.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasIPAccessList.spec
<sup><sup>[↩ Parent](#atlasipaccesslist)</sup></sup>



AtlasIPAccessListSpec defines the desired state of AtlasIPAccessList.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasipaccesslistspecentriesindex">entries</a></b></td>
        <td>[]object</td>
        <td>
          Entries is the list of IP Access to be managed.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasipaccesslistspecconnectionsecret">connectionSecret</a></b></td>
        <td>object</td>
        <td>
          Name of the secret containing Atlas API private and public keys.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasipaccesslistspecexternalprojectref">externalProjectRef</a></b></td>
        <td>object</td>
        <td>
          externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasipaccesslistspecprojectref">projectRef</a></b></td>
        <td>object</td>
        <td>
          projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasIPAccessList.spec.entries[index]
<sup><sup>[↩ Parent](#atlasipaccesslistspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>awsSecurityGroup</b></td>
        <td>string</td>
        <td>
          Unique identifier of AWS security group in this access list entry.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>cidrBlock</b></td>
        <td>string</td>
        <td>
          Range of IP addresses in CIDR notation in this access list entry.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>comment</b></td>
        <td>string</td>
        <td>
          Comment associated with this access list entry.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deleteAfterDate</b></td>
        <td>string</td>
        <td>
          Date and time after which Atlas deletes the temporary access list entry.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ipAddress</b></td>
        <td>string</td>
        <td>
          Entry using an IP address in this access list entry.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasIPAccessList.spec.connectionSecret
<sup><sup>[↩ Parent](#atlasipaccesslistspec)</sup></sup>



Name of the secret containing Atlas API private and public keys.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasIPAccessList.spec.externalProjectRef
<sup><sup>[↩ Parent](#atlasipaccesslistspec)</sup></sup>



externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID is the Atlas project ID.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasIPAccessList.spec.projectRef
<sup><sup>[↩ Parent](#atlasipaccesslistspec)</sup></sup>



projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasIPAccessList.status
<sup><sup>[↩ Parent](#atlasipaccesslist)</sup></sup>



AtlasIPAccessListStatus is the most recent observed status of the AtlasIPAccessList cluster. Read-only.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasipaccessliststatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasipaccessliststatusentriesindex">entries</a></b></td>
        <td>[]object</td>
        <td>
          Status is the state of the ip access list<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasIPAccessList.status.conditions[index]
<sup><sup>[↩ Parent](#atlasipaccessliststatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasIPAccessList.status.entries[index]
<sup><sup>[↩ Parent](#atlasipaccessliststatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>entry</b></td>
        <td>string</td>
        <td>
          Entry is the ip access Atlas is managing<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status is the correspondent state of the entry<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>

## AtlasNetworkContainer
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasNetworkContainer is the Schema for the AtlasNetworkContainer API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasNetworkContainer</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkcontainerspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasNetworkContainerSpec defines the desired state of an AtlasNetworkContainer.<br/>
          <br/>
            <i>Validations</i>:<li>(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef)): must define only one project reference through externalProjectRef or projectRef</li><li>(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef): must define a local connection secret when referencing an external project</li><li>(self.provider == 'GCP' && !has(self.region)) || (self.provider != 'GCP'): must not set region for GCP containers</li><li>((self.provider == 'AWS' || self.provider == 'AZURE') && has(self.region)) || (self.provider == 'GCP'): must set region for AWS and Azure containers</li><li>(self.id == oldSelf.id) || (!has(self.id) && !has(oldSelf.id)): id is immutable</li><li>(self.region == oldSelf.region) || (!has(self.region) && !has(oldSelf.region)): region is immutable</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkcontainerstatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasNetworkContainerStatus is a status for the AtlasNetworkContainer Custom resource.
Not the one included in the AtlasProject<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkContainer.spec
<sup><sup>[↩ Parent](#atlasnetworkcontainer)</sup></sup>



AtlasNetworkContainerSpec defines the desired state of an AtlasNetworkContainer.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>provider</b></td>
        <td>enum</td>
        <td>
          Provider is the name of the cloud provider hosting the network container.<br/>
          <br/>
            <i>Enum</i>: AWS, GCP, AZURE<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>cidrBlock</b></td>
        <td>string</td>
        <td>
          Atlas CIDR. It needs to be set if ContainerID is not set.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkcontainerspecconnectionsecret">connectionSecret</a></b></td>
        <td>object</td>
        <td>
          Name of the secret containing Atlas API private and public keys.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkcontainerspecexternalprojectref">externalProjectRef</a></b></td>
        <td>object</td>
        <td>
          externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID is the container identifier for an already existent network container to be managed by the operator.
This field can be used in conjunction with cidrBlock to update the cidrBlock of an existing container.
This field is immutable.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkcontainerspecprojectref">projectRef</a></b></td>
        <td>object</td>
        <td>
          projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          ContainerRegion is the provider region name of Atlas network peer container in Atlas region format
This is required by AWS and Azure, but not used by GCP.
This field is immutable, Atlas does not admit network container changes.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkContainer.spec.connectionSecret
<sup><sup>[↩ Parent](#atlasnetworkcontainerspec)</sup></sup>



Name of the secret containing Atlas API private and public keys.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasNetworkContainer.spec.externalProjectRef
<sup><sup>[↩ Parent](#atlasnetworkcontainerspec)</sup></sup>



externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID is the Atlas project ID.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasNetworkContainer.spec.projectRef
<sup><sup>[↩ Parent](#atlasnetworkcontainerspec)</sup></sup>



projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkContainer.status
<sup><sup>[↩ Parent](#atlasnetworkcontainer)</sup></sup>



AtlasNetworkContainerStatus is a status for the AtlasNetworkContainer Custom resource.
Not the one included in the AtlasProject

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasnetworkcontainerstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID record the identifier of the container in Atlas<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>provisioned</b></td>
        <td>boolean</td>
        <td>
          Provisioned is true when clusters have been deployed to the container before
the last reconciliation<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkContainer.status.conditions[index]
<sup><sup>[↩ Parent](#atlasnetworkcontainerstatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasNetworkPeering
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasNetworkPeering is the Schema for the AtlasNetworkPeering API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasNetworkPeering</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkpeeringspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasNetworkPeeringSpec defines the desired state of AtlasNetworkPeering<br/>
          <br/>
            <i>Validations</i>:<li>(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef)): must define only one project reference through externalProjectRef or projectRef</li><li>(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef): must define a local connection secret when referencing an external project</li><li>(has(self.containerRef.name) && !has(self.containerRef.id)) || (!has(self.containerRef.name) && has(self.containerRef.id)): must either have a container Atlas id or Kubernetes name, but not both (or neither)</li><li>(self.containerRef.name == oldSelf.containerRef.name) || (!has(self.containerRef.name) && !has(oldSelf.containerRef.name)): container ref name is immutable</li><li>(self.containerRef.id == oldSelf.containerRef.id) || (!has(self.containerRef.id) && !has(oldSelf.containerRef.id)): container ref id is immutable</li><li>(self.id == oldSelf.id) || (!has(self.id) && !has(oldSelf.id)): id is immutable</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkpeeringstatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasNetworkPeeringStatus is a status for the AtlasNetworkPeering Custom resource.
Not the one included in the AtlasProject<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.spec
<sup><sup>[↩ Parent](#atlasnetworkpeering)</sup></sup>



AtlasNetworkPeeringSpec defines the desired state of AtlasNetworkPeering

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasnetworkpeeringspeccontainerref">containerRef</a></b></td>
        <td>object</td>
        <td>
          ContainerDualReference refers to a Network Container either by Kubernetes name or Atlas ID.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>provider</b></td>
        <td>enum</td>
        <td>
          Name of the cloud service provider for which you want to create the network peering service.<br/>
          <br/>
            <i>Enum</i>: AWS, GCP, AZURE<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkpeeringspecawsconfiguration">awsConfiguration</a></b></td>
        <td>object</td>
        <td>
          AWSConfiguration is the specific AWS settings for network peering.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkpeeringspecazureconfiguration">azureConfiguration</a></b></td>
        <td>object</td>
        <td>
          AzureConfiguration is the specific Azure settings for network peering.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkpeeringspecconnectionsecret">connectionSecret</a></b></td>
        <td>object</td>
        <td>
          Name of the secret containing Atlas API private and public keys.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkpeeringspecexternalprojectref">externalProjectRef</a></b></td>
        <td>object</td>
        <td>
          externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkpeeringspecgcpconfiguration">gcpConfiguration</a></b></td>
        <td>object</td>
        <td>
          GCPConfiguration is the specific Google Cloud settings for network peering.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID is the peering identifier for an already existent network peering to be managed by the operator.
This field is immutable.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkpeeringspecprojectref">projectRef</a></b></td>
        <td>object</td>
        <td>
          projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.spec.containerRef
<sup><sup>[↩ Parent](#atlasnetworkpeeringspec)</sup></sup>



ContainerDualReference refers to a Network Container either by Kubernetes name or Atlas ID.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID is the Atlas identifier of the Network Container Atlas resource this Peering Connection relies on.
Use either name or ID, not both.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the container Kubernetes resource, must be present in the same namespace.
Use either name or ID, not both.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.spec.awsConfiguration
<sup><sup>[↩ Parent](#atlasnetworkpeeringspec)</sup></sup>



AWSConfiguration is the specific AWS settings for network peering.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>accepterRegionName</b></td>
        <td>string</td>
        <td>
          AccepterRegionName is the provider region name of user's vpc in AWS native region format.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>awsAccountId</b></td>
        <td>string</td>
        <td>
          AccountID of the user's vpc.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>routeTableCidrBlock</b></td>
        <td>string</td>
        <td>
          User VPC CIDR.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>vpcId</b></td>
        <td>string</td>
        <td>
          AWS VPC ID.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.spec.azureConfiguration
<sup><sup>[↩ Parent](#atlasnetworkpeeringspec)</sup></sup>



AzureConfiguration is the specific Azure settings for network peering.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>azureDirectoryId</b></td>
        <td>string</td>
        <td>
          AzureDirectoryID is the unique identifier for an Azure AD directory.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>azureSubscriptionId</b></td>
        <td>string</td>
        <td>
          AzureSubscriptionID is the unique identifier of the Azure subscription in which the VNet resides.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>resourceGroupName</b></td>
        <td>string</td>
        <td>
          ResourceGroupName is the name of your Azure resource group.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>vNetName</b></td>
        <td>string</td>
        <td>
          VNetName is name of your Azure VNet. Its applicable only for Azure.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.spec.connectionSecret
<sup><sup>[↩ Parent](#atlasnetworkpeeringspec)</sup></sup>



Name of the secret containing Atlas API private and public keys.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.spec.externalProjectRef
<sup><sup>[↩ Parent](#atlasnetworkpeeringspec)</sup></sup>



externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID is the Atlas project ID.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.spec.gcpConfiguration
<sup><sup>[↩ Parent](#atlasnetworkpeeringspec)</sup></sup>



GCPConfiguration is the specific Google Cloud settings for network peering.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>gcpProjectId</b></td>
        <td>string</td>
        <td>
          User GCP Project ID. Its applicable only for GCP.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>networkName</b></td>
        <td>string</td>
        <td>
          GCP Network Peer Name. Its applicable only for GCP.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.spec.projectRef
<sup><sup>[↩ Parent](#atlasnetworkpeeringspec)</sup></sup>



projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.status
<sup><sup>[↩ Parent](#atlasnetworkpeering)</sup></sup>



AtlasNetworkPeeringStatus is a status for the AtlasNetworkPeering Custom resource.
Not the one included in the AtlasProject

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasnetworkpeeringstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkpeeringstatusawsstatus">awsStatus</a></b></td>
        <td>object</td>
        <td>
          AWSStatus contains AWS only related status information<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkpeeringstatusazurestatus">azureStatus</a></b></td>
        <td>object</td>
        <td>
          AzureStatus contains Azure only related status information<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasnetworkpeeringstatusgcpstatus">gcpStatus</a></b></td>
        <td>object</td>
        <td>
          GCPStatus contains GCP only related status information<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID recrods the identified of the peer created by Atlas<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status describes the last status seen for the network peering setup<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.status.conditions[index]
<sup><sup>[↩ Parent](#atlasnetworkpeeringstatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.status.awsStatus
<sup><sup>[↩ Parent](#atlasnetworkpeeringstatus)</sup></sup>



AWSStatus contains AWS only related status information

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>connectionId</b></td>
        <td>string</td>
        <td>
          ConnectionID is the AWS VPC peering connection ID<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>vpcId</b></td>
        <td>string</td>
        <td>
          VpcID is AWS VPC id on the Atlas side<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.status.azureStatus
<sup><sup>[↩ Parent](#atlasnetworkpeeringstatus)</sup></sup>



AzureStatus contains Azure only related status information

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>azureSubscriptionIDpcId</b></td>
        <td>string</td>
        <td>
          AzureSubscriptionID is Azure Subscription id on the Atlas side<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>vNetName</b></td>
        <td>string</td>
        <td>
          VnetName is Azure network on the Atlas side<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasNetworkPeering.status.gcpStatus
<sup><sup>[↩ Parent](#atlasnetworkpeeringstatus)</sup></sup>



GCPStatus contains GCP only related status information

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>gcpProjectID</b></td>
        <td>string</td>
        <td>
          GCPProjectID is GCP project on the Atlas side<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>networkName</b></td>
        <td>string</td>
        <td>
          NetworkName is GCP network on the Atlas side<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasOrgSettings
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>








<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasOrgSettings</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasorgsettingsspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasOrgSettingsSpec defines the desired state of AtlasOrgSettings.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasorgsettingsstatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasOrgSettingsStatus defines the observed state of AtlasOrgSettings.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasOrgSettings.spec
<sup><sup>[↩ Parent](#atlasorgsettings)</sup></sup>



AtlasOrgSettingsSpec defines the desired state of AtlasOrgSettings.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>orgID</b></td>
        <td>string</td>
        <td>
          OrgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiAccessListRequired</b></td>
        <td>boolean</td>
        <td>
          ApiAccessListRequired Flag that indicates whether to require API operations to originate from an IP Address added to the API access list for the specified organization.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasorgsettingsspecconnectionsecretref">connectionSecretRef</a></b></td>
        <td>object</td>
        <td>
          ConnectionSecretRef is the name of the Kubernetes Secret which contains the information about the way to connect to Atlas (Public & Private API keys).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>genAIFeaturesEnabled</b></td>
        <td>boolean</td>
        <td>
          GenAIFeaturesEnabled Flag that indicates whether this organization has access to generative AI features. This setting only applies to Atlas Commercial and is enabled by default.
Once this setting is turned on, Project Owners may be able to enable or disable individual AI features at the project level.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxServiceAccountSecretValidityInHours</b></td>
        <td>integer</td>
        <td>
          MaxServiceAccountSecretValidityInHours Number that represents the maximum period before expiry in hours for new Atlas Admin API Service Account secrets within the specified organization.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>multiFactorAuthRequired</b></td>
        <td>boolean</td>
        <td>
          MultiFactorAuthRequired Flag that indicates whether to require users to set up Multi-Factor Authentication (MFA) before accessing the specified organization.
To learn more, see: https://www.mongodb.com/docs/atlas/security-multi-factor-authentication/.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>restrictEmployeeAccess</b></td>
        <td>boolean</td>
        <td>
          RestrictEmployeeAccess Flag that indicates whether to block MongoDB Support from accessing Atlas infrastructure and cluster logs for any deployment in the specified organization without explicit permission.
Once this setting is turned on, you can grant MongoDB Support a 24-hour bypass access to the Atlas deployment to resolve support issues.
To learn more, see: https://www.mongodb.com/docs/atlas/security-restrict-support-access/.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>securityContact</b></td>
        <td>string</td>
        <td>
          SecurityContact String that specifies a single email address for the specified organization to receive security-related notifications.
Specifying a security contact does not grant them authorization or access to Atlas for security decisions or approvals.
An empty string is valid and clears the existing security contact (if any).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>streamsCrossGroupEnabled</b></td>
        <td>boolean</td>
        <td>
          StreamsCrossGroupEnabled Flag that indicates whether a group's Atlas Stream Processing instances in this organization can create connections to other group's clusters in the same organization.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasOrgSettings.spec.connectionSecretRef
<sup><sup>[↩ Parent](#atlasorgsettingsspec)</sup></sup>



ConnectionSecretRef is the name of the Kubernetes Secret which contains the information about the way to connect to Atlas (Public & Private API keys).

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasOrgSettings.status
<sup><sup>[↩ Parent](#atlasorgsettings)</sup></sup>



AtlasOrgSettingsStatus defines the observed state of AtlasOrgSettings.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasorgsettingsstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions holding the status details<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasOrgSettings.status.conditions[index]
<sup><sup>[↩ Parent](#atlasorgsettingsstatus)</sup></sup>



Condition contains details for one aspect of the current state of this API Resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          lastTransitionTime is the last time the condition transitioned from one status to another.
This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          message is a human readable message indicating details about the transition.
This may be an empty string.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          reason contains a programmatic identifier indicating the reason for the condition's last transition.
Producers of specific condition types may define expected values and meanings for this field,
and whether the values are considered a guaranteed API.
The value should be a CamelCase string.
This field may not be empty.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>enum</td>
        <td>
          status of the condition, one of True, False, Unknown.<br/>
          <br/>
            <i>Enum</i>: True, False, Unknown<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type of condition in CamelCase or in foo.example.com/CamelCase.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          observedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Minimum</i>: 0<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasPrivateEndpoint
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






The AtlasPrivateEndpoint custom resource definition (CRD) defines a desired [Private Endpoint](https://www.mongodb.com/docs/atlas/security-private-endpoint/#std-label-private-endpoint-overview) configuration for an Atlas project.
It allows a private connection between your cloud provider and Atlas that doesn't send information through a public network.

You can use private endpoints to create a unidirectional connection to Atlas clusters from your virtual network.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasPrivateEndpoint</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasprivateendpointspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasPrivateEndpointSpec is the specification of the desired configuration of a project private endpoint<br/>
          <br/>
            <i>Validations</i>:<li>(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef)): must define only one project reference through externalProjectRef or projectRef</li><li>(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef): must define a local connection secret when referencing an external project</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprivateendpointstatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasPrivateEndpointStatus is the most recent observed status of the AtlasPrivateEndpoint cluster. Read-only.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.spec
<sup><sup>[↩ Parent](#atlasprivateendpoint)</sup></sup>



AtlasPrivateEndpointSpec is the specification of the desired configuration of a project private endpoint

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>provider</b></td>
        <td>enum</td>
        <td>
          Name of the cloud service provider for which you want to create the private endpoint service.<br/>
          <br/>
            <i>Enum</i>: AWS, GCP, AZURE<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          Region of the chosen cloud provider in which you want to create the private endpoint service.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasprivateendpointspecawsconfigurationindex">awsConfiguration</a></b></td>
        <td>[]object</td>
        <td>
          AWSConfiguration is the specific AWS settings for the private endpoint.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprivateendpointspecazureconfigurationindex">azureConfiguration</a></b></td>
        <td>[]object</td>
        <td>
          AzureConfiguration is the specific Azure settings for the private endpoint.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprivateendpointspecconnectionsecret">connectionSecret</a></b></td>
        <td>object</td>
        <td>
          Name of the secret containing Atlas API private and public keys.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprivateendpointspecexternalprojectref">externalProjectRef</a></b></td>
        <td>object</td>
        <td>
          externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprivateendpointspecgcpconfigurationindex">gcpConfiguration</a></b></td>
        <td>[]object</td>
        <td>
          GCPConfiguration is the specific Google Cloud settings for the private endpoint.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprivateendpointspecprojectref">projectRef</a></b></td>
        <td>object</td>
        <td>
          projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.spec.awsConfiguration[index]
<sup><sup>[↩ Parent](#atlasprivateendpointspec)</sup></sup>



AWSPrivateEndpointConfiguration holds the AWS configuration done on customer network.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID that identifies the private endpoint's network interface that someone added to this private endpoint service.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.spec.azureConfiguration[index]
<sup><sup>[↩ Parent](#atlasprivateendpointspec)</sup></sup>



AzurePrivateEndpointConfiguration holds the Azure configuration done on customer network.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID that identifies the private endpoint's network interface that someone added to this private endpoint service.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>ipAddress</b></td>
        <td>string</td>
        <td>
          IP address of the private endpoint in your Azure VNet that someone added to this private endpoint service.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.spec.connectionSecret
<sup><sup>[↩ Parent](#atlasprivateendpointspec)</sup></sup>



Name of the secret containing Atlas API private and public keys.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.spec.externalProjectRef
<sup><sup>[↩ Parent](#atlasprivateendpointspec)</sup></sup>



externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID is the Atlas project ID.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.spec.gcpConfiguration[index]
<sup><sup>[↩ Parent](#atlasprivateendpointspec)</sup></sup>



GCPPrivateEndpointConfiguration holds the GCP configuration done on customer network.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasprivateendpointspecgcpconfigurationindexendpointsindex">endpoints</a></b></td>
        <td>[]object</td>
        <td>
          Endpoints is the list of individual private endpoints that comprise this endpoint group.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>groupName</b></td>
        <td>string</td>
        <td>
          GroupName is the label that identifies a set of endpoints.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>projectId</b></td>
        <td>string</td>
        <td>
          ProjectID that identifies the Google Cloud project in which you created the endpoints.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.spec.gcpConfiguration[index].endpoints[index]
<sup><sup>[↩ Parent](#atlasprivateendpointspecgcpconfigurationindex)</sup></sup>



GCPPrivateEndpoint holds the GCP forwarding rules configured on customer network.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>ipAddress</b></td>
        <td>string</td>
        <td>
          IP address to which this Google Cloud consumer forwarding rule resolves.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name that identifies the Google Cloud consumer forwarding rule that you created.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.spec.projectRef
<sup><sup>[↩ Parent](#atlasprivateendpointspec)</sup></sup>



projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.status
<sup><sup>[↩ Parent](#atlasprivateendpoint)</sup></sup>



AtlasPrivateEndpointStatus is the most recent observed status of the AtlasPrivateEndpoint cluster. Read-only.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasprivateendpointstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasprivateendpointstatusendpointsindex">endpoints</a></b></td>
        <td>[]object</td>
        <td>
          Endpoints are the status of the endpoints connected to the service<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>error</b></td>
        <td>string</td>
        <td>
          Error is the description of the failure occurred when configuring the private endpoint<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resourceId</b></td>
        <td>string</td>
        <td>
          ResourceID is the root-relative path that identifies of the Atlas Azure Private Link Service<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceAttachmentNames</b></td>
        <td>[]string</td>
        <td>
          ServiceAttachmentNames is the list of URLs that identifies endpoints that Atlas can use to access one service across the private connection<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceId</b></td>
        <td>string</td>
        <td>
          ServiceID is the unique identifier of the private endpoint service in Atlas<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceName</b></td>
        <td>string</td>
        <td>
          ServiceName is the unique identifier of the Amazon Web Services (AWS) PrivateLink endpoint service or Azure Private Link Service managed by Atlas<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceStatus</b></td>
        <td>string</td>
        <td>
          ServiceStatus is the state of the private endpoint service<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.status.conditions[index]
<sup><sup>[↩ Parent](#atlasprivateendpointstatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.status.endpoints[index]
<sup><sup>[↩ Parent](#atlasprivateendpointstatus)</sup></sup>



EndpointInterfaceStatus is the most recent observed status the interfaces attached to the configured service. Read-only.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>ID</b></td>
        <td>string</td>
        <td>
          ID is the external identifier set on the specification to configure the interface<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>InterfaceStatus</b></td>
        <td>string</td>
        <td>
          InterfaceStatus is the state of the private endpoint interface<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>connectionName</b></td>
        <td>string</td>
        <td>
          ConnectionName is the label that Atlas generates that identifies the Azure private endpoint connection<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>error</b></td>
        <td>string</td>
        <td>
          Error is the description of the failure occurred when configuring the private endpoint<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprivateendpointstatusendpointsindexgcpforwardingrulesindex">gcpForwardingRules</a></b></td>
        <td>[]object</td>
        <td>
          GCPForwardingRules is the status of the customer GCP private endpoint(forwarding rules)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasPrivateEndpoint.status.endpoints[index].gcpForwardingRules[index]
<sup><sup>[↩ Parent](#atlasprivateendpointstatusendpointsindex)</sup></sup>



GCPForwardingRule is the most recent observed status the GCP forwarding rules configured for an interface. Read-only.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the Google Cloud consumer forwarding rule that you created.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          State of the MongoDB Atlas endpoint group.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasProject
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasProject is the Schema for the atlasprojects API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasProject</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasProjectSpec defines the desired state of Project in Atlas<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasProjectStatus defines the observed state of AtlasProject<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec
<sup><sup>[↩ Parent](#atlasproject)</sup></sup>



AtlasProjectSpec defines the desired state of Project in Atlas

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of the Project that is created in Atlas by the Operator if it doesn't exist yet.
The name length must not exceed 64 characters. The name must contain only letters, numbers, spaces, dashes, and underscores.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: Name cannot be modified after project creation</li>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>alertConfigurationSyncEnabled</b></td>
        <td>boolean</td>
        <td>
          AlertConfigurationSyncEnabled is a flag that enables/disables Alert Configurations sync for the current Project.
If true - project alert configurations will be synced according to AlertConfigurations.
If not - alert configurations will not be modified by the operator. They can be managed through the API, CLI, and UI.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecalertconfigurationsindex">alertConfigurations</a></b></td>
        <td>[]object</td>
        <td>
          AlertConfiguration is a list of Alert Configurations configured for the current Project.
If you use this setting, you must also set spec.alertConfigurationSyncEnabled to true for Atlas Kubernetes
Operator to modify project alert configurations.
If you omit or leave this setting empty, Atlas Kubernetes Operator doesn't alter the project's alert
configurations. If creating a project, Atlas applies the default project alert configurations.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecauditing">auditing</a></b></td>
        <td>object</td>
        <td>
          Auditing represents MongoDB Maintenance Windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecbackupcompliancepolicyref">backupCompliancePolicyRef</a></b></td>
        <td>object</td>
        <td>
          BackupCompliancePolicyRef is a reference to the backup compliance custom resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspeccloudprovideraccessrolesindex">cloudProviderAccessRoles</a></b></td>
        <td>[]object</td>
        <td>
          CloudProviderAccessRoles is a list of Cloud Provider Access Roles configured for the current Project.
Deprecated: This configuration was deprecated in favor of CloudProviderIntegrations<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspeccloudproviderintegrationsindex">cloudProviderIntegrations</a></b></td>
        <td>[]object</td>
        <td>
          CloudProviderIntegrations is a list of Cloud Provider Integration configured for the current Project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecconnectionsecretref">connectionSecretRef</a></b></td>
        <td>object</td>
        <td>
          ConnectionSecret is the name of the Kubernetes Secret which contains the information about the way to connect to
Atlas (organization ID, API keys). The default Operator connection configuration will be used if not provided.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspeccustomrolesindex">customRoles</a></b></td>
        <td>[]object</td>
        <td>
          CustomRoles lets you create and change custom roles in your cluster.
Use custom roles to specify custom sets of actions that the Atlas built-in roles can't describe.
Deprecated: Migrate to the AtlasCustomRoles custom resource in accordance with the migration guide
at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecencryptionatrest">encryptionAtRest</a></b></td>
        <td>object</td>
        <td>
          EncryptionAtRest allows to set encryption for AWS, Azure and GCP providers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecintegrationsindex">integrations</a></b></td>
        <td>[]object</td>
        <td>
          Integrations is a list of MongoDB Atlas integrations for the project.
Deprecated: Migrate to the AtlasThirdPartyIntegration custom resource in accordance with the migration guide
at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecmaintenancewindow">maintenanceWindow</a></b></td>
        <td>object</td>
        <td>
          MaintenanceWindow allows to specify a preferred time in the week to run maintenance operations. See more
information at https://www.mongodb.com/docs/atlas/reference/api/maintenance-windows/<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecnetworkpeersindex">networkPeers</a></b></td>
        <td>[]object</td>
        <td>
          NetworkPeers is a list of Network Peers configured for the current Project.
Deprecated: Migrate to the AtlasNetworkPeering and AtlasNetworkContainer custom resources in accordance with
the migration guide at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecprivateendpointsindex">privateEndpoints</a></b></td>
        <td>[]object</td>
        <td>
          PrivateEndpoints is a list of Private Endpoints configured for the current Project.
Deprecated: Migrate to the AtlasPrivateEndpoint Custom Resource in accordance with the migration guide
at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecprojectipaccesslistindex">projectIpAccessList</a></b></td>
        <td>[]object</td>
        <td>
          ProjectIPAccessList allows the use of the IP Access List for a Project. See more information at
https://docs.atlas.mongodb.com/reference/api/ip-access-list/add-entries-to-access-list/
Deprecated: Migrate to the AtlasIPAccessList Custom Resource in accordance with the migration guide
at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>regionUsageRestrictions</b></td>
        <td>enum</td>
        <td>
          RegionUsageRestrictions designate the project's AWS region when using Atlas for Government.
This parameter should not be used with commercial Atlas.
In Atlas for Government, not setting this field (defaulting to NONE) means the project is restricted to COMMERCIAL_FEDRAMP_REGIONS_ONLY.<br/>
          <br/>
            <i>Enum</i>: NONE, GOV_REGIONS_ONLY, COMMERCIAL_FEDRAMP_REGIONS_ONLY<br/>
            <i>Default</i>: NONE<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecsettings">settings</a></b></td>
        <td>object</td>
        <td>
          Settings allows the configuration of the Project Settings.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecteamsindex">teams</a></b></td>
        <td>[]object</td>
        <td>
          Teams enable you to grant project access roles to multiple users.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>withDefaultAlertsSettings</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether Atlas Kubernetes Operator creates a project with the default alert configurations.
If you use this setting, you must also set spec.alertConfigurationSyncEnabled to true for Atlas Kubernetes
Operator to modify project alert configurations.
If you set this parameter to false when you create a project, Atlas doesn't add the default alert configurations
to your project.
This setting has no effect on existing projects.<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecx509certref">x509CertRef</a></b></td>
        <td>object</td>
        <td>
          X509CertRef is a reference to the Kubernetes Secret which contains PEM-encoded CA certificate.
Atlas Kubernetes Operator watches secrets only with the label atlas.mongodb.com/type=credentials to avoid
watching unnecessary secrets.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.alertConfigurations[index]
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          If omitted, the configuration is disabled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>eventTypeName</b></td>
        <td>string</td>
        <td>
          The type of event that will trigger an alert.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecalertconfigurationsindexmatchersindex">matchers</a></b></td>
        <td>[]object</td>
        <td>
          You can filter using the matchers array only when the EventTypeName specifies an event for a host, replica set, or sharded cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecalertconfigurationsindexmetricthreshold">metricThreshold</a></b></td>
        <td>object</td>
        <td>
          MetricThreshold  causes an alert to be triggered.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecalertconfigurationsindexnotificationsindex">notifications</a></b></td>
        <td>[]object</td>
        <td>
          Notifications are sending when an alert condition is detected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>severityOverride</b></td>
        <td>enum</td>
        <td>
          SeverityOverride optionally overrides the default severity level for an alert.<br/>
          <br/>
            <i>Enum</i>: INFO, WARNING, ERROR, CRITICAL<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecalertconfigurationsindexthreshold">threshold</a></b></td>
        <td>object</td>
        <td>
          Threshold  causes an alert to be triggered.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.alertConfigurations[index].matchers[index]
<sup><sup>[↩ Parent](#atlasprojectspecalertconfigurationsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fieldName</b></td>
        <td>string</td>
        <td>
          Name of the field in the target object to match on.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          The operator to test the field’s value.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Value to test with the specified operator.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.alertConfigurations[index].metricThreshold
<sup><sup>[↩ Parent](#atlasprojectspecalertconfigurationsindex)</sup></sup>



MetricThreshold  causes an alert to be triggered.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>threshold</b></td>
        <td>string</td>
        <td>
          Threshold value outside which an alert will be triggered.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>metricName</b></td>
        <td>string</td>
        <td>
          Name of the metric to check.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>string</td>
        <td>
          This must be set to AVERAGE. Atlas computes the current metric value as an average.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Operator to apply when checking the current metric value against the threshold value.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>units</b></td>
        <td>string</td>
        <td>
          The units for the threshold value.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.alertConfigurations[index].notifications[index]
<sup><sup>[↩ Parent](#atlasprojectspecalertconfigurationsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasprojectspecalertconfigurationsindexnotificationsindexapitokenref">apiTokenRef</a></b></td>
        <td>object</td>
        <td>
          Secret containing a Slack API token or Bot token. Populated for the SLACK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>channelName</b></td>
        <td>string</td>
        <td>
          Slack channel name. Populated for the SLACK notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecalertconfigurationsindexnotificationsindexdatadogapikeyref">datadogAPIKeyRef</a></b></td>
        <td>object</td>
        <td>
          Secret containing a Datadog API Key. Found in the Datadog dashboard. Populated for the DATADOG notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>datadogRegion</b></td>
        <td>string</td>
        <td>
          Region that indicates which API URL to use.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>delayMin</b></td>
        <td>integer</td>
        <td>
          Number of minutes to wait after an alert condition is detected before sending out the first notification.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>emailAddress</b></td>
        <td>string</td>
        <td>
          Email address to which alert notifications are sent. Populated for the EMAIL notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>emailEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag indicating if email notifications should be sent. Populated for ORG, GROUP, and USER notifications types.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>flowName</b></td>
        <td>string</td>
        <td>
          Flowdock flow name in lower-case letters.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecalertconfigurationsindexnotificationsindexflowdockapitokenref">flowdockApiTokenRef</a></b></td>
        <td>object</td>
        <td>
          The Flowdock personal API token. Populated for the FLOWDOCK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>intervalMin</b></td>
        <td>integer</td>
        <td>
          Number of minutes to wait between successive notifications for unacknowledged alerts that are not resolved.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mobileNumber</b></td>
        <td>string</td>
        <td>
          Mobile number to which alert notifications are sent. Populated for the SMS notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecalertconfigurationsindexnotificationsindexopsgenieapikeyref">opsGenieApiKeyRef</a></b></td>
        <td>object</td>
        <td>
          OpsGenie API Key. Populated for the OPS_GENIE notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>opsGenieRegion</b></td>
        <td>string</td>
        <td>
          Region that indicates which API URL to use.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>orgName</b></td>
        <td>string</td>
        <td>
          Flowdock organization name in lower-case letters. This is the name that appears after www.flowdock.com/app/ in the URL string. Populated for the FLOWDOCK notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>roles</b></td>
        <td>[]string</td>
        <td>
          The following roles grant privileges within a project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecalertconfigurationsindexnotificationsindexservicekeyref">serviceKeyRef</a></b></td>
        <td>object</td>
        <td>
          PagerDuty service key. Populated for the PAGER_DUTY notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>smsEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag indicating if text message notifications should be sent. Populated for ORG, GROUP, and USER notifications types.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>teamId</b></td>
        <td>string</td>
        <td>
          Unique identifier of a team.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>teamName</b></td>
        <td>string</td>
        <td>
          Label for the team that receives this notification.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>typeName</b></td>
        <td>string</td>
        <td>
          Type of alert notification.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>username</b></td>
        <td>string</td>
        <td>
          Name of the Atlas user to which to send notifications. Only a user in the project that owns the alert configuration is allowed here. Populated for the USER notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecalertconfigurationsindexnotificationsindexvictoropssecretref">victorOpsSecretRef</a></b></td>
        <td>object</td>
        <td>
          Secret containing a VictorOps API key and Routing key. Populated for the VICTOR_OPS notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.alertConfigurations[index].notifications[index].apiTokenRef
<sup><sup>[↩ Parent](#atlasprojectspecalertconfigurationsindexnotificationsindex)</sup></sup>



Secret containing a Slack API token or Bot token. Populated for the SLACK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.alertConfigurations[index].notifications[index].datadogAPIKeyRef
<sup><sup>[↩ Parent](#atlasprojectspecalertconfigurationsindexnotificationsindex)</sup></sup>



Secret containing a Datadog API Key. Found in the Datadog dashboard. Populated for the DATADOG notifications type.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.alertConfigurations[index].notifications[index].flowdockApiTokenRef
<sup><sup>[↩ Parent](#atlasprojectspecalertconfigurationsindexnotificationsindex)</sup></sup>



The Flowdock personal API token. Populated for the FLOWDOCK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.alertConfigurations[index].notifications[index].opsGenieApiKeyRef
<sup><sup>[↩ Parent](#atlasprojectspecalertconfigurationsindexnotificationsindex)</sup></sup>



OpsGenie API Key. Populated for the OPS_GENIE notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.alertConfigurations[index].notifications[index].serviceKeyRef
<sup><sup>[↩ Parent](#atlasprojectspecalertconfigurationsindexnotificationsindex)</sup></sup>



PagerDuty service key. Populated for the PAGER_DUTY notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.alertConfigurations[index].notifications[index].victorOpsSecretRef
<sup><sup>[↩ Parent](#atlasprojectspecalertconfigurationsindexnotificationsindex)</sup></sup>



Secret containing a VictorOps API key and Routing key. Populated for the VICTOR_OPS notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.alertConfigurations[index].threshold
<sup><sup>[↩ Parent](#atlasprojectspecalertconfigurationsindex)</sup></sup>



Threshold  causes an alert to be triggered.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Operator to apply when checking the current metric value against the threshold value.
It accepts the following values: GREATER_THAN, LESS_THAN.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>threshold</b></td>
        <td>string</td>
        <td>
          Threshold value outside which an alert will be triggered.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>units</b></td>
        <td>string</td>
        <td>
          The units for the threshold value.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.auditing
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



Auditing represents MongoDB Maintenance Windows.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>auditAuthorizationSuccess</b></td>
        <td>boolean</td>
        <td>
          Indicates whether the auditing system captures successful authentication attempts for audit filters using the "atype" : "authCheck" auditing event.
For more information, see auditAuthorizationSuccess.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>auditFilter</b></td>
        <td>string</td>
        <td>
          JSON-formatted audit filter used by the project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Denotes whether the project associated with the {GROUP-ID} has database auditing enabled.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.backupCompliancePolicyRef
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



BackupCompliancePolicyRef is a reference to the backup compliance custom resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.cloudProviderAccessRoles[index]
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



CloudProviderAccessRole define an integration to a cloud provider
DEPRECATED: This type is deprecated in favor of CloudProviderIntegration

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>providerName</b></td>
        <td>string</td>
        <td>
          ProviderName is the name of the cloud provider. Currently only AWS is supported.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>iamAssumedRoleArn</b></td>
        <td>string</td>
        <td>
          IamAssumedRoleArn is the ARN of the IAM role that is assumed by the Atlas cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.cloudProviderIntegrations[index]
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



CloudProviderIntegration define an integration to a cloud provider

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>providerName</b></td>
        <td>string</td>
        <td>
          ProviderName is the name of the cloud provider. Currently only AWS is supported.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>iamAssumedRoleArn</b></td>
        <td>string</td>
        <td>
          IamAssumedRoleArn is the ARN of the IAM role that is assumed by the Atlas cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.connectionSecretRef
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



ConnectionSecret is the name of the Kubernetes Secret which contains the information about the way to connect to
Atlas (organization ID, API keys). The default Operator connection configuration will be used if not provided.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.customRoles[index]
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



CustomRole lets you create and change a custom role in your cluster.
Use custom roles to specify custom sets of actions that the Atlas built-in roles can't describe.
Deprecated: Migrate to the AtlasCustomRoles custom resource in accordance with the migration guide
at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the role. This name must be unique for this custom role in this project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspeccustomrolesindexactionsindex">actions</a></b></td>
        <td>[]object</td>
        <td>
          List of the individual privilege actions that the role grants.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspeccustomrolesindexinheritedrolesindex">inheritedRoles</a></b></td>
        <td>[]object</td>
        <td>
          List of the built-in roles that this custom role inherits.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.customRoles[index].actions[index]
<sup><sup>[↩ Parent](#atlasprojectspeccustomrolesindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the privilege action.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspeccustomrolesindexactionsindexresourcesindex">resources</a></b></td>
        <td>[]object</td>
        <td>
          List of resources on which you grant the action.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasProject.spec.customRoles[index].actions[index].resources[index]
<sup><sup>[↩ Parent](#atlasprojectspeccustomrolesindexactionsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>cluster</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to grant the action on the cluster resource. If true, MongoDB Cloud ignores Database and Collection parameters.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>collection</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the collection on which you grant the action to one MongoDB user.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the database on which you grant the action to one MongoDB user.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.customRoles[index].inheritedRoles[index]
<sup><sup>[↩ Parent](#atlasprojectspeccustomrolesindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>database</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the database on which someone grants the action to one MongoDB user.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the role inherited.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasProject.spec.encryptionAtRest
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



EncryptionAtRest allows to set encryption for AWS, Azure and GCP providers.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasprojectspecencryptionatrestawskms">awsKms</a></b></td>
        <td>object</td>
        <td>
          AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecencryptionatrestazurekeyvault">azureKeyVault</a></b></td>
        <td>object</td>
        <td>
          AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecencryptionatrestgooglecloudkms">googleCloudKms</a></b></td>
        <td>object</td>
        <td>
          GoogleCloudKms specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.encryptionAtRest.awsKms
<sup><sup>[↩ Parent](#atlasprojectspecencryptionatrest)</sup></sup>



AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Specifies whether Encryption at Rest is enabled for an Atlas project.
To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          The AWS region in which the AWS customer master key exists.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecencryptionatrestawskmssecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          A reference to as Secret containing the AccessKeyID, SecretAccessKey, CustomerMasterKeyID and RoleID fields<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>valid</b></td>
        <td>boolean</td>
        <td>
          Specifies whether the encryption key set for the provider is valid and may be used to encrypt and decrypt data.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.encryptionAtRest.awsKms.secretRef
<sup><sup>[↩ Parent](#atlasprojectspecencryptionatrestawskms)</sup></sup>



A reference to as Secret containing the AccessKeyID, SecretAccessKey, CustomerMasterKeyID and RoleID fields

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.encryptionAtRest.azureKeyVault
<sup><sup>[↩ Parent](#atlasprojectspecencryptionatrest)</sup></sup>



AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>azureEnvironment</b></td>
        <td>string</td>
        <td>
          The Azure environment where the Azure account credentials reside. Valid values are the following: AZURE, AZURE_CHINA, AZURE_GERMANY<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>clientID</b></td>
        <td>string</td>
        <td>
          The Client ID, also known as the application ID, for an Azure application associated with the Azure AD tenant.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Specifies whether Encryption at Rest is enabled for an Atlas project.
To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resourceGroupName</b></td>
        <td>string</td>
        <td>
          The name of the Azure Resource group that contains an Azure Key Vault.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecencryptionatrestazurekeyvaultsecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          A reference to as Secret containing the SubscriptionID, KeyVaultName, KeyIdentifier, Secret fields<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tenantID</b></td>
        <td>string</td>
        <td>
          The unique identifier for an Azure AD tenant within an Azure subscription.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.encryptionAtRest.azureKeyVault.secretRef
<sup><sup>[↩ Parent](#atlasprojectspecencryptionatrestazurekeyvault)</sup></sup>



A reference to as Secret containing the SubscriptionID, KeyVaultName, KeyIdentifier, Secret fields

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.encryptionAtRest.googleCloudKms
<sup><sup>[↩ Parent](#atlasprojectspecencryptionatrest)</sup></sup>



GoogleCloudKms specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Specifies whether Encryption at Rest is enabled for an Atlas project.
To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecencryptionatrestgooglecloudkmssecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          A reference to as Secret containing the ServiceAccountKey, KeyVersionResourceID fields<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.encryptionAtRest.googleCloudKms.secretRef
<sup><sup>[↩ Parent](#atlasprojectspecencryptionatrestgooglecloudkms)</sup></sup>



A reference to as Secret containing the ServiceAccountKey, KeyVersionResourceID fields

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.integrations[index]
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



Integration for the project between Atlas and a third party service.
Deprecated: Migrate to the AtlasThirdPartyIntegration custom resource in accordance with the migration guide
at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>accountId</b></td>
        <td>string</td>
        <td>
          Unique 40-hexadecimal digit string that identifies your New Relic account.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecintegrationsindexapikeyref">apiKeyRef</a></b></td>
        <td>object</td>
        <td>
          Reference to a Kubernetes Secret containing your API Key for Datadog, OpsGenie or Victor Ops.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecintegrationsindexapitokenref">apiTokenRef</a></b></td>
        <td>object</td>
        <td>
          Reference to a Kubernetes Secret containing the Key that allows Atlas to access your Slack account.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>channelName</b></td>
        <td>string</td>
        <td>
          Name of the Slack channel to which Atlas sends alert notifications.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether someone has activated the Prometheus integration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>flowName</b></td>
        <td>string</td>
        <td>
          DEPRECATED: Flowdock flow name.
This field has been removed from Atlas, and has no effect.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecintegrationsindexlicensekeyref">licenseKeyRef</a></b></td>
        <td>object</td>
        <td>
          Reference to a Kubernetes Secret containing your Unique 40-hexadecimal digit string that identifies your New Relic license.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>microsoftTeamsWebhookUrl</b></td>
        <td>string</td>
        <td>
          Endpoint web address of the Microsoft Teams webhook to which Atlas sends notifications.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>orgName</b></td>
        <td>string</td>
        <td>
          DEPRECATED: Flowdock organization name.
This field has been removed from Atlas, and has no effect.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecintegrationsindexpasswordref">passwordRef</a></b></td>
        <td>object</td>
        <td>
          Reference to a Kubernetes Secret containing the password to allow Atlas to access your Prometheus account.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecintegrationsindexreadtokenref">readTokenRef</a></b></td>
        <td>object</td>
        <td>
          Reference to a Kubernetes Secret containing the query key associated with your New Relic account.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          Region code indicating which regional API Atlas uses to access PagerDuty, Datadog, or OpsGenie.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecintegrationsindexroutingkeyref">routingKeyRef</a></b></td>
        <td>object</td>
        <td>
          Reference to a Kubernetes Secret containing the Routing key associated with your Splunk On-Call account.
Used for Victor Ops.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecintegrationsindexsecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          Reference to a Kubernetes Secret containing the secret for your Webhook.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceDiscovery</b></td>
        <td>string</td>
        <td>
          Desired method to discover the Prometheus service.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecintegrationsindexservicekeyref">serviceKeyRef</a></b></td>
        <td>object</td>
        <td>
          Reference to a Kubernetes Secret containing the service key associated with your PagerDuty account.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>teamName</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies your Slack team.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Third Party Integration type such as Slack, New Relic, etc.
Each integration type requires a distinct set of configuration fields.
For example, if you set type to DATADOG, you must configure only datadog subfields.<br/>
          <br/>
            <i>Enum</i>: PAGER_DUTY, SLACK, DATADOG, NEW_RELIC, OPS_GENIE, VICTOR_OPS, FLOWDOCK, WEBHOOK, MICROSOFT_TEAMS, PROMETHEUS<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>url</b></td>
        <td>string</td>
        <td>
          Endpoint web address to which Atlas sends notifications.
Used for Webhooks.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>username</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies your Prometheus incoming webhook.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecintegrationsindexwritetokenref">writeTokenRef</a></b></td>
        <td>object</td>
        <td>
          Reference to a Kubernetes Secret containing the insert key associated with your New Relic account.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.integrations[index].apiKeyRef
<sup><sup>[↩ Parent](#atlasprojectspecintegrationsindex)</sup></sup>



Reference to a Kubernetes Secret containing your API Key for Datadog, OpsGenie or Victor Ops.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.integrations[index].apiTokenRef
<sup><sup>[↩ Parent](#atlasprojectspecintegrationsindex)</sup></sup>



Reference to a Kubernetes Secret containing the Key that allows Atlas to access your Slack account.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.integrations[index].licenseKeyRef
<sup><sup>[↩ Parent](#atlasprojectspecintegrationsindex)</sup></sup>



Reference to a Kubernetes Secret containing your Unique 40-hexadecimal digit string that identifies your New Relic license.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.integrations[index].passwordRef
<sup><sup>[↩ Parent](#atlasprojectspecintegrationsindex)</sup></sup>



Reference to a Kubernetes Secret containing the password to allow Atlas to access your Prometheus account.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.integrations[index].readTokenRef
<sup><sup>[↩ Parent](#atlasprojectspecintegrationsindex)</sup></sup>



Reference to a Kubernetes Secret containing the query key associated with your New Relic account.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.integrations[index].routingKeyRef
<sup><sup>[↩ Parent](#atlasprojectspecintegrationsindex)</sup></sup>



Reference to a Kubernetes Secret containing the Routing key associated with your Splunk On-Call account.
Used for Victor Ops.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.integrations[index].secretRef
<sup><sup>[↩ Parent](#atlasprojectspecintegrationsindex)</sup></sup>



Reference to a Kubernetes Secret containing the secret for your Webhook.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.integrations[index].serviceKeyRef
<sup><sup>[↩ Parent](#atlasprojectspecintegrationsindex)</sup></sup>



Reference to a Kubernetes Secret containing the service key associated with your PagerDuty account.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.integrations[index].writeTokenRef
<sup><sup>[↩ Parent](#atlasprojectspecintegrationsindex)</sup></sup>



Reference to a Kubernetes Secret containing the insert key associated with your New Relic account.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.maintenanceWindow
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



MaintenanceWindow allows to specify a preferred time in the week to run maintenance operations. See more
information at https://www.mongodb.com/docs/atlas/reference/api/maintenance-windows/

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>autoDefer</b></td>
        <td>boolean</td>
        <td>
          Flag indicating whether any scheduled project maintenance should be deferred automatically for one week.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>dayOfWeek</b></td>
        <td>integer</td>
        <td>
          Day of the week when you would like the maintenance window to start as a 1-based integer.
Sunday 1, Monday 2, Tuesday 3, Wednesday 4, Thursday 5, Friday 6, Saturday 7.<br/>
          <br/>
            <i>Minimum</i>: 1<br/>
            <i>Maximum</i>: 7<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>defer</b></td>
        <td>boolean</td>
        <td>
          Flag indicating whether the next scheduled project maintenance should be deferred for one week.
Cannot be specified if startASAP is true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hourOfDay</b></td>
        <td>integer</td>
        <td>
          Hour of the day when you would like the maintenance window to start.
This parameter uses the 24-hour clock, where midnight is 0, noon is 12.<br/>
          <br/>
            <i>Minimum</i>: 0<br/>
            <i>Maximum</i>: 23<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>startASAP</b></td>
        <td>boolean</td>
        <td>
          Flag indicating whether project maintenance has been directed to start immediately.
Cannot be specified if defer is true<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.networkPeers[index]
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



NetworkPeer configured for the current Project.
Deprecated: Migrate to the AtlasNetworkPeering and AtlasNetworkContainer custom resources in accordance with
the migration guide at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>accepterRegionName</b></td>
        <td>string</td>
        <td>
          AccepterRegionName is the provider region name of user's VPC.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>atlasCidrBlock</b></td>
        <td>string</td>
        <td>
          Atlas CIDR. It needs to be set if ContainerID is not set.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>awsAccountId</b></td>
        <td>string</td>
        <td>
          AccountID of the user's VPC.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>azureDirectoryId</b></td>
        <td>string</td>
        <td>
          AzureDirectoryID is the unique identifier for an Azure AD directory.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>azureSubscriptionId</b></td>
        <td>string</td>
        <td>
          AzureSubscriptionID is the unique identifier of the Azure subscription in which the VNet resides.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>containerId</b></td>
        <td>string</td>
        <td>
          ID of the network peer container. If not set, operator will create a new container with ContainerRegion and AtlasCIDRBlock input.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>containerRegion</b></td>
        <td>string</td>
        <td>
          ContainerRegion is the provider region name of Atlas network peer container. If not set, AccepterRegionName is used.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gcpProjectId</b></td>
        <td>string</td>
        <td>
          User GCP Project ID. Its applicable only for GCP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>networkName</b></td>
        <td>string</td>
        <td>
          GCP Network Peer Name. Its applicable only for GCP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>providerName</b></td>
        <td>string</td>
        <td>
          ProviderName is the name of the provider. If not set, it will be set to "AWS".<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resourceGroupName</b></td>
        <td>string</td>
        <td>
          ResourceGroupName is the name of your Azure resource group.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>routeTableCidrBlock</b></td>
        <td>string</td>
        <td>
          User VPC CIDR.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>vnetName</b></td>
        <td>string</td>
        <td>
          VNetName is name of your Azure VNet. Its applicable only for Azure.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>vpcId</b></td>
        <td>string</td>
        <td>
          AWS VPC ID.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.privateEndpoints[index]
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



PrivateEndpoint is a list of Private Endpoints configured for the current Project.
Deprecated: Migrate to the AtlasPrivateEndpoint Custom Resource in accordance with the migration guide
at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>provider</b></td>
        <td>enum</td>
        <td>
          Cloud provider for which you want to retrieve a private endpoint service. Atlas accepts AWS, GCP, or AZURE.<br/>
          <br/>
            <i>Enum</i>: AWS, GCP, AZURE, TENANT<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          Cloud provider region for which you want to create the private endpoint service.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>endpointGroupName</b></td>
        <td>string</td>
        <td>
          Unique identifier of the endpoint group. The endpoint group encompasses all the endpoints that you created in Google Cloud.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecprivateendpointsindexendpointsindex">endpoints</a></b></td>
        <td>[]object</td>
        <td>
          Collection of individual private endpoints that comprise your endpoint group.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gcpProjectId</b></td>
        <td>string</td>
        <td>
          Unique identifier of the Google Cloud project in which you created your endpoints.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Unique identifier of the private endpoint you created in your AWS VPC or Azure VNet.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ip</b></td>
        <td>string</td>
        <td>
          Private IP address of the private endpoint network interface you created in your Azure VNet.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.privateEndpoints[index].endpoints[index]
<sup><sup>[↩ Parent](#atlasprojectspecprivateendpointsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>endpointName</b></td>
        <td>string</td>
        <td>
          Forwarding rule that corresponds to the endpoint you created in Google Cloud.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ipAddress</b></td>
        <td>string</td>
        <td>
          Private IP address of the endpoint you created in Google Cloud.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.projectIpAccessList[index]
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



IPAccessList allows the use of the IP Access List for a Project. See more information at
https://docs.atlas.mongodb.com/reference/api/ip-access-list/add-entries-to-access-list/
Deprecated: Migrate to the AtlasIPAccessList Custom Resource in accordance with the migration guide
at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>awsSecurityGroup</b></td>
        <td>string</td>
        <td>
          Unique identifier of AWS security group in this access list entry.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>cidrBlock</b></td>
        <td>string</td>
        <td>
          Range of IP addresses in CIDR notation in this access list entry.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>comment</b></td>
        <td>string</td>
        <td>
          Comment associated with this access list entry.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deleteAfterDate</b></td>
        <td>string</td>
        <td>
          Timestamp in ISO 8601 date and time format in UTC after which Atlas deletes the temporary access list entry.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ipAddress</b></td>
        <td>string</td>
        <td>
          Entry using an IP address in this access list entry.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.settings
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



Settings allows the configuration of the Project Settings.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>isCollectDatabaseSpecificsStatisticsEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to collect database-specific metrics for the specified project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isDataExplorerEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to enable the Data Explorer for the specified project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isExtendedStorageSizesEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to enable extended storage sizes for the specified project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isPerformanceAdvisorEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to enable the Performance Advisor and Profiler for the specified project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isRealtimePerformancePanelEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to enable the Real Time Performance Panel for the specified project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isSchemaAdvisorEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to enable the Schema Advisor for the specified project.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.teams[index]
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>roles</b></td>
        <td>[]enum</td>
        <td>
          Roles the users in the team has within the project.<br/>
          <br/>
            <i>Enum</i>: GROUP_OWNER, GROUP_CLUSTER_MANAGER, GROUP_DATA_ACCESS_ADMIN, GROUP_DATA_ACCESS_READ_WRITE, GROUP_DATA_ACCESS_READ_ONLY, GROUP_READ_ONLY<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasprojectspecteamsindexteamref">teamRef</a></b></td>
        <td>object</td>
        <td>
          Reference to the AtlasTeam custom resource which will be assigned to the project.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasProject.spec.teams[index].teamRef
<sup><sup>[↩ Parent](#atlasprojectspecteamsindex)</sup></sup>



Reference to the AtlasTeam custom resource which will be assigned to the project.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.spec.x509CertRef
<sup><sup>[↩ Parent](#atlasprojectspec)</sup></sup>



X509CertRef is a reference to the Kubernetes Secret which contains PEM-encoded CA certificate.
Atlas Kubernetes Operator watches secrets only with the label atlas.mongodb.com/type=credentials to avoid
watching unnecessary secrets.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status
<sup><sup>[↩ Parent](#atlasproject)</sup></sup>



AtlasProjectStatus defines the observed state of AtlasProject

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasprojectstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusalertconfigurationsindex">alertConfigurations</a></b></td>
        <td>[]object</td>
        <td>
          AlertConfigurations contains a list of alert configuration statuses<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>authModes</b></td>
        <td>[]string</td>
        <td>
          AuthModes contains a list of configured authentication modes
"SCRAM" is default authentication method and requires a password for each user
"X509" signifies that self-managed X.509 authentication is configured<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatuscloudproviderintegrationsindex">cloudProviderIntegrations</a></b></td>
        <td>[]object</td>
        <td>
          CloudProviderIntegrations contains a list of configured cloud provider access roles. AWS support only<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatuscustomrolesindex">customRoles</a></b></td>
        <td>[]object</td>
        <td>
          CustomRoles contains a list of custom roles statuses<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusexpiredipaccesslistindex">expiredIpAccessList</a></b></td>
        <td>[]object</td>
        <td>
          The list of IP Access List entries that are expired due to 'deleteAfterDate' being less than the current date.
Note, that this field is updated by the Atlas Operator only after specification changes<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the Atlas Project<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusnetworkpeersindex">networkPeers</a></b></td>
        <td>[]object</td>
        <td>
          The list of network peers that are configured for current project<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusprivateendpointsindex">privateEndpoints</a></b></td>
        <td>[]object</td>
        <td>
          The list of private endpoints configured for current project<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusprometheus">prometheus</a></b></td>
        <td>object</td>
        <td>
          Prometheus contains the status for Prometheus integration
including the prometheusDiscoveryURL<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusteamsindex">teams</a></b></td>
        <td>[]object</td>
        <td>
          Teams contains a list of teams assignment statuses<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.conditions[index]
<sup><sup>[↩ Parent](#atlasprojectstatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.alertConfigurations[index]
<sup><sup>[↩ Parent](#atlasprojectstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>acknowledgedUntil</b></td>
        <td>string</td>
        <td>
          The date through which the alert has been acknowledged. Will not be present if the alert has never been acknowledged.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>acknowledgementComment</b></td>
        <td>string</td>
        <td>
          The comment left by the user who acknowledged the alert. Will not be present if the alert has never been acknowledged.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>acknowledgingUsername</b></td>
        <td>string</td>
        <td>
          The username of the user who acknowledged the alert. Will not be present if the alert has never been acknowledged.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>alertConfigId</b></td>
        <td>string</td>
        <td>
          ID of the alert configuration that triggered this alert.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>clusterId</b></td>
        <td>string</td>
        <td>
          The ID of the cluster to which this alert applies. Only present for alerts of type BACKUP, REPLICA_SET, and CLUSTER.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>clusterName</b></td>
        <td>string</td>
        <td>
          The name the cluster to which this alert applies. Only present for alerts of type BACKUP, REPLICA_SET, and CLUSTER.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>created</b></td>
        <td>string</td>
        <td>
          Timestamp in ISO 8601 date and time format in UTC when this alert configuration was created.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusalertconfigurationsindexcurrentvalue">currentValue</a></b></td>
        <td>object</td>
        <td>
          CurrentValue represents current value of the metric that triggered the alert. Only present for alerts of type HOST_METRIC.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          If omitted, the configuration is disabled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>errorMessage</b></td>
        <td>string</td>
        <td>
          ErrorMessage is massage if the alert configuration is in an incorrect state.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>eventTypeName</b></td>
        <td>string</td>
        <td>
          The type of event that will trigger an alert.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>groupId</b></td>
        <td>string</td>
        <td>
          Unique identifier of the project that owns this alert configuration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostId</b></td>
        <td>string</td>
        <td>
          ID of the host to which the metric pertains. Only present for alerts of type HOST, HOST_METRIC, and REPLICA_SET.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostnameAndPort</b></td>
        <td>string</td>
        <td>
          The hostname and port of each host to which the alert applies. Only present for alerts of type HOST, HOST_METRIC, and REPLICA_SET.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Unique identifier.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastNotified</b></td>
        <td>string</td>
        <td>
          When the last notification was sent for this alert. Only present if notifications have been sent.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusalertconfigurationsindexmatchersindex">matchers</a></b></td>
        <td>[]object</td>
        <td>
          You can filter using the matchers array only when the EventTypeName specifies an event for a host, replica set, or sharded cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>metricName</b></td>
        <td>string</td>
        <td>
          The name of the measurement whose value went outside the threshold. Only present if eventTypeName is set to OUTSIDE_METRIC_THRESHOLD.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusalertconfigurationsindexmetricthreshold">metricThreshold</a></b></td>
        <td>object</td>
        <td>
          MetricThreshold  causes an alert to be triggered.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusalertconfigurationsindexnotificationsindex">notifications</a></b></td>
        <td>[]object</td>
        <td>
          Notifications are sending when an alert condition is detected.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>replicaSetName</b></td>
        <td>string</td>
        <td>
          Name of the replica set. Only present for alerts of type HOST, HOST_METRIC, BACKUP, and REPLICA_SET.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resolved</b></td>
        <td>string</td>
        <td>
          When the alert was closed. Only present if the status is CLOSED.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>severityOverride</b></td>
        <td>string</td>
        <td>
          Severity of the alert.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sourceTypeName</b></td>
        <td>string</td>
        <td>
          For alerts of the type BACKUP, the type of server being backed up.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          The current state of the alert. Possible values are: TRACKING, OPEN, CLOSED, CANCELED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusalertconfigurationsindexthreshold">threshold</a></b></td>
        <td>object</td>
        <td>
          Threshold  causes an alert to be triggered.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>updated</b></td>
        <td>string</td>
        <td>
          Timestamp in ISO 8601 date and time format in UTC when this alert configuration was last updated.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.alertConfigurations[index].currentValue
<sup><sup>[↩ Parent](#atlasprojectstatusalertconfigurationsindex)</sup></sup>



CurrentValue represents current value of the metric that triggered the alert. Only present for alerts of type HOST_METRIC.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>number</b></td>
        <td>string</td>
        <td>
          The value of the metric.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>units</b></td>
        <td>string</td>
        <td>
          The units for the value. Depends on the type of metric.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.alertConfigurations[index].matchers[index]
<sup><sup>[↩ Parent](#atlasprojectstatusalertconfigurationsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fieldName</b></td>
        <td>string</td>
        <td>
          Name of the field in the target object to match on.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          The operator to test the field’s value.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Value to test with the specified operator.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.alertConfigurations[index].metricThreshold
<sup><sup>[↩ Parent](#atlasprojectstatusalertconfigurationsindex)</sup></sup>



MetricThreshold  causes an alert to be triggered.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>threshold</b></td>
        <td>string</td>
        <td>
          Threshold value outside which an alert will be triggered.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>metricName</b></td>
        <td>string</td>
        <td>
          Name of the metric to check.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>string</td>
        <td>
          This must be set to AVERAGE. Atlas computes the current metric value as an average.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Operator to apply when checking the current metric value against the threshold value.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>units</b></td>
        <td>string</td>
        <td>
          The units for the threshold value.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.alertConfigurations[index].notifications[index]
<sup><sup>[↩ Parent](#atlasprojectstatusalertconfigurationsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>apiToken</b></td>
        <td>string</td>
        <td>
          Slack API token or Bot token. Populated for the SLACK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>channelName</b></td>
        <td>string</td>
        <td>
          Slack channel name. Populated for the SLACK notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>datadogApiKey</b></td>
        <td>string</td>
        <td>
          Datadog API Key. Found in the Datadog dashboard. Populated for the DATADOG notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>datadogRegion</b></td>
        <td>string</td>
        <td>
          Region that indicates which API URL to use<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>delayMin</b></td>
        <td>integer</td>
        <td>
          Number of minutes to wait after an alert condition is detected before sending out the first notification.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>emailAddress</b></td>
        <td>string</td>
        <td>
          Email address to which alert notifications are sent. Populated for the EMAIL notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>emailEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag indicating if email notifications should be sent. Populated for ORG, GROUP, and USER notifications types.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>flowName</b></td>
        <td>string</td>
        <td>
          Flowdock flow namse in lower-case letters.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>flowdockApiToken</b></td>
        <td>string</td>
        <td>
          The Flowdock personal API token. Populated for the FLOWDOCK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>intervalMin</b></td>
        <td>integer</td>
        <td>
          Number of minutes to wait between successive notifications for unacknowledged alerts that are not resolved.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mobileNumber</b></td>
        <td>string</td>
        <td>
          Mobile number to which alert notifications are sent. Populated for the SMS notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>opsGenieApiKey</b></td>
        <td>string</td>
        <td>
          Opsgenie API Key. Populated for the OPS_GENIE notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>opsGenieRegion</b></td>
        <td>string</td>
        <td>
          Region that indicates which API URL to use.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>orgName</b></td>
        <td>string</td>
        <td>
          Flowdock organization name in lower-case letters. This is the name that appears after www.flowdock.com/app/ in the URL string. Populated for the FLOWDOCK notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>roles</b></td>
        <td>[]string</td>
        <td>
          The following roles grant privileges within a project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceKey</b></td>
        <td>string</td>
        <td>
          PagerDuty service key. Populated for the PAGER_DUTY notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>smsEnabled</b></td>
        <td>boolean</td>
        <td>
          Flag indicating if text message notifications should be sent. Populated for ORG, GROUP, and USER notifications types.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>teamId</b></td>
        <td>string</td>
        <td>
          Unique identifier of a team.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>teamName</b></td>
        <td>string</td>
        <td>
          Label for the team that receives this notification.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>typeName</b></td>
        <td>string</td>
        <td>
          Type of alert notification.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>username</b></td>
        <td>string</td>
        <td>
          Name of the Atlas user to which to send notifications. Only a user in the project that owns the alert configuration is allowed here. Populated for the USER notifications type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>victorOpsApiKey</b></td>
        <td>string</td>
        <td>
          VictorOps API key. Populated for the VICTOR_OPS notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>victorOpsRoutingKey</b></td>
        <td>string</td>
        <td>
          VictorOps routing key. Populated for the VICTOR_OPS notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.alertConfigurations[index].threshold
<sup><sup>[↩ Parent](#atlasprojectstatusalertconfigurationsindex)</sup></sup>



Threshold  causes an alert to be triggered.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Operator to apply when checking the current metric value against the threshold value. it accepts the following values: GREATER_THAN, LESS_THAN<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>threshold</b></td>
        <td>string</td>
        <td>
          Threshold value outside which an alert will be triggered.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>units</b></td>
        <td>string</td>
        <td>
          The units for the threshold value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.cloudProviderIntegrations[index]
<sup><sup>[↩ Parent](#atlasprojectstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>atlasAssumedRoleExternalId</b></td>
        <td>string</td>
        <td>
          Unique external ID that MongoDB Atlas uses when it assumes the IAM role in your Amazon Web Services account.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>providerName</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the cloud provider of the role.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>atlasAWSAccountArn</b></td>
        <td>string</td>
        <td>
          Amazon Resource Name that identifies the Amazon Web Services user account that MongoDB Atlas uses when it assumes the Identity and Access Management role.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>authorizedDate</b></td>
        <td>string</td>
        <td>
          Date and time when someone authorized this role for the specified cloud service provider. This parameter expresses its value in the ISO 8601 timestamp format in UTC.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>createdDate</b></td>
        <td>string</td>
        <td>
          Date and time when someone created this role for the specified cloud service provider. This parameter expresses its value in the ISO 8601 timestamp format in UTC.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>errorMessage</b></td>
        <td>string</td>
        <td>
          Application error message returned.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatuscloudproviderintegrationsindexfeatureusagesindex">featureUsages</a></b></td>
        <td>[]object</td>
        <td>
          List that contains application features associated with this Amazon Web Services Identity and Access Management role.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>iamAssumedRoleArn</b></td>
        <td>string</td>
        <td>
          Amazon Resource Name that identifies the Amazon Web Services Identity and Access Management role that MongoDB Cloud assumes when it accesses resources in your AWS account.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>roleId</b></td>
        <td>string</td>
        <td>
          Unique 24-hexadecimal digit string that identifies the role.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Provision status of the service account.
Values are IN_PROGRESS, COMPLETE, FAILED, or NOT_INITIATED.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.cloudProviderIntegrations[index].featureUsages[index]
<sup><sup>[↩ Parent](#atlasprojectstatuscloudproviderintegrationsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>featureId</b></td>
        <td>string</td>
        <td>
          Identifying characteristics about the data lake linked to this Amazon Web Services Identity and Access Management role.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>featureType</b></td>
        <td>string</td>
        <td>
          Human-readable label that describes one MongoDB Cloud feature linked to this Amazon Web Services Identity and Access Management role.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.customRoles[index]
<sup><sup>[↩ Parent](#atlasprojectstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Role name which is unique<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          The status of the given custom role (OK or FAILED)<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>error</b></td>
        <td>string</td>
        <td>
          The message when the custom role is in the FAILED status<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.expiredIpAccessList[index]
<sup><sup>[↩ Parent](#atlasprojectstatus)</sup></sup>



IPAccessList allows the use of the IP Access List for a Project. See more information at
https://docs.atlas.mongodb.com/reference/api/ip-access-list/add-entries-to-access-list/
Deprecated: Migrate to the AtlasIPAccessList Custom Resource in accordance with the migration guide
at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>awsSecurityGroup</b></td>
        <td>string</td>
        <td>
          Unique identifier of AWS security group in this access list entry.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>cidrBlock</b></td>
        <td>string</td>
        <td>
          Range of IP addresses in CIDR notation in this access list entry.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>comment</b></td>
        <td>string</td>
        <td>
          Comment associated with this access list entry.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deleteAfterDate</b></td>
        <td>string</td>
        <td>
          Timestamp in ISO 8601 date and time format in UTC after which Atlas deletes the temporary access list entry.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ipAddress</b></td>
        <td>string</td>
        <td>
          Entry using an IP address in this access list entry.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.networkPeers[index]
<sup><sup>[↩ Parent](#atlasprojectstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Unique identifier for NetworkPeer.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>providerName</b></td>
        <td>string</td>
        <td>
          Cloud provider for which you want to retrieve a network peer.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          Region for which you want to create the network peer. It isn't needed for GCP<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>atlasGcpProjectId</b></td>
        <td>string</td>
        <td>
          ProjectID of Atlas container. Applicable only for GCP. It's needed to add network peer connection.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>atlasNetworkName</b></td>
        <td>string</td>
        <td>
          Atlas Network Name. Applicable only for GCP. It's needed to add network peer connection.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>connectionId</b></td>
        <td>string</td>
        <td>
          Unique identifier of the network peer connection. Applicable only for AWS.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>containerId</b></td>
        <td>string</td>
        <td>
          ContainerID of Atlas network peer container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>errorMessage</b></td>
        <td>string</td>
        <td>
          Error state of the network peer. Applicable only for GCP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>errorState</b></td>
        <td>string</td>
        <td>
          Error state of the network peer. Applicable only for Azure.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>errorStateName</b></td>
        <td>string</td>
        <td>
          Error state of the network peer. Applicable only for AWS.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gcpProjectId</b></td>
        <td>string</td>
        <td>
          ProjectID of the user's vpc. Applicable only for GCP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the network peer. Applicable only for GCP and Azure.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>statusName</b></td>
        <td>string</td>
        <td>
          Status of the network peer. Applicable only for AWS.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>vpc</b></td>
        <td>string</td>
        <td>
          VPC is general purpose field for storing the name of the VPC.
VPC is vpcID for AWS, user networkName for GCP, and vnetName for Azure.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.privateEndpoints[index]
<sup><sup>[↩ Parent](#atlasprojectstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>provider</b></td>
        <td>string</td>
        <td>
          Cloud provider for which you want to retrieve a private endpoint service. Atlas accepts AWS or AZURE.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          Cloud provider region for which you want to create the private endpoint service.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasprojectstatusprivateendpointsindexendpointsindex">endpoints</a></b></td>
        <td>[]object</td>
        <td>
          Collection of individual GCP private endpoints that comprise your network endpoint group.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Unique identifier for AWS or AZURE Private Link Connection.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>interfaceEndpointId</b></td>
        <td>string</td>
        <td>
          Unique identifier of the AWS or Azure Private Link Interface Endpoint.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceAttachmentNames</b></td>
        <td>[]string</td>
        <td>
          Unique alphanumeric and special character strings that identify the service attachments associated with the GCP Private Service Connect endpoint service.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceName</b></td>
        <td>string</td>
        <td>
          Name of the AWS or Azure Private Link Service that Atlas manages.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceResourceId</b></td>
        <td>string</td>
        <td>
          Unique identifier of the Azure Private Link Service (for AWS the same as ID).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.privateEndpoints[index].endpoints[index]
<sup><sup>[↩ Parent](#atlasprojectstatusprivateendpointsindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>endpointName</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the Google Cloud consumer forwarding rule that you created.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>ipAddress</b></td>
        <td>string</td>
        <td>
          One Private Internet Protocol version 4 (IPv4) address to which this Google Cloud consumer forwarding rule resolves.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          State of the MongoDB Atlas endpoint group when MongoDB Cloud received this request.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasProject.status.prometheus
<sup><sup>[↩ Parent](#atlasprojectstatus)</sup></sup>



Prometheus contains the status for Prometheus integration
including the prometheusDiscoveryURL

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>prometheusDiscoveryURL</b></td>
        <td>string</td>
        <td>
          URL from which Prometheus fetches the targets.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Protocol scheme used for Prometheus requests.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.teams[index]
<sup><sup>[↩ Parent](#atlasprojectstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasprojectstatusteamsindexteamref">teamRef</a></b></td>
        <td>object</td>
        <td>
          ResourceRefNamespaced is a reference to a Kubernetes Resource that allows to configure the namespace<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasProject.status.teams[index].teamRef
<sup><sup>[↩ Parent](#atlasprojectstatusteamsindex)</sup></sup>



ResourceRefNamespaced is a reference to a Kubernetes Resource that allows to configure the namespace

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasSearchIndexConfig
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasSearchIndexConfig is the Schema for the AtlasSearchIndexConfig API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasSearchIndexConfig</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlassearchindexconfigspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasSearchIndexConfigSpec defines the desired state of AtlasSearchIndexConfig.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlassearchindexconfigstatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasSearchIndexConfigStatus defines the observed state of AtlasSearchIndexConfig.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasSearchIndexConfig.spec
<sup><sup>[↩ Parent](#atlassearchindexconfig)</sup></sup>



AtlasSearchIndexConfigSpec defines the desired state of AtlasSearchIndexConfig.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>analyzer</b></td>
        <td>enum</td>
        <td>
          Specific pre-defined method chosen to convert database field text into searchable words. This conversion reduces the text of fields into the smallest units of text.
These units are called a term or token. This process, known as tokenization, involves a variety of changes made to the text in fields:
- extracting words
- removing punctuation
- removing accents
- hanging to lowercase
- removing common words
- reducing words to their root form (stemming)
- changing words to their base form (lemmatization) MongoDB Cloud uses the selected process to build the Atlas Search index<br/>
          <br/>
            <i>Enum</i>: lucene.standard, lucene.simple, lucene.whitespace, lucene.keyword, lucene.arabic, lucene.armenian, lucene.basque, lucene.bengali, lucene.brazilian, lucene.bulgarian, lucene.catalan, lucene.chinese, lucene.cjk, lucene.czech, lucene.danish, lucene.dutch, lucene.english, lucene.finnish, lucene.french, lucene.galician, lucene.german, lucene.greek, lucene.hindi, lucene.hungarian, lucene.indonesian, lucene.irish, lucene.italian, lucene.japanese, lucene.korean, lucene.kuromoji, lucene.latvian, lucene.lithuanian, lucene.morfologik, lucene.nori, lucene.norwegian, lucene.persian, lucene.portuguese, lucene.romanian, lucene.russian, lucene.smartcn, lucene.sorani, lucene.spanish, lucene.swedish, lucene.thai, lucene.turkish, lucene.ukrainian<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlassearchindexconfigspecanalyzersindex">analyzers</a></b></td>
        <td>[]object</td>
        <td>
          List of user-defined methods to convert database field text into searchable words.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>searchAnalyzer</b></td>
        <td>enum</td>
        <td>
          Method applied to identify words when searching this index.<br/>
          <br/>
            <i>Enum</i>: lucene.standard, lucene.simple, lucene.whitespace, lucene.keyword, lucene.arabic, lucene.armenian, lucene.basque, lucene.bengali, lucene.brazilian, lucene.bulgarian, lucene.catalan, lucene.chinese, lucene.cjk, lucene.czech, lucene.danish, lucene.dutch, lucene.english, lucene.finnish, lucene.french, lucene.galician, lucene.german, lucene.greek, lucene.hindi, lucene.hungarian, lucene.indonesian, lucene.irish, lucene.italian, lucene.japanese, lucene.korean, lucene.kuromoji, lucene.latvian, lucene.lithuanian, lucene.morfologik, lucene.nori, lucene.norwegian, lucene.persian, lucene.portuguese, lucene.romanian, lucene.russian, lucene.smartcn, lucene.sorani, lucene.spanish, lucene.swedish, lucene.thai, lucene.turkish, lucene.ukrainian<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storedSource</b></td>
        <td>JSON</td>
        <td>
          Flag that indicates whether to store all fields (true) on Atlas Search. By default, Atlas doesn't store (false) the fields on Atlas Search.
Alternatively, you can specify an object that only contains the list of fields to store (include) or not store (exclude) on Atlas Search.
To learn more, see documentation: https://www.mongodb.com/docs/atlas/atlas-search/stored-source-definition/<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasSearchIndexConfig.spec.analyzers[index]
<sup><sup>[↩ Parent](#atlassearchindexconfigspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable name that identifies the custom analyzer. Names must be unique within an index, and must not start with any of the following strings:
"lucene.", "builtin.", "mongodb."<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlassearchindexconfigspecanalyzersindextokenizer">tokenizer</a></b></td>
        <td>object</td>
        <td>
          Tokenizer that you want to use to create tokens. Tokens determine how Atlas Search splits up text into discrete chunks for indexing.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>charFilters</b></td>
        <td>JSON</td>
        <td>
          Filters that examine text one character at a time and perform filtering operations.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tokenFilters</b></td>
        <td>JSON</td>
        <td>
          Filter that performs operations such as:
- Stemming, which reduces related words, such as "talking", "talked", and "talks" to their root word "talk".
- Redaction, the removal of sensitive information from public documents<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasSearchIndexConfig.spec.analyzers[index].tokenizer
<sup><sup>[↩ Parent](#atlassearchindexconfigspecanalyzersindex)</sup></sup>



Tokenizer that you want to use to create tokens. Tokens determine how Atlas Search splits up text into discrete chunks for indexing.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Human-readable label that identifies this tokenizer type.<br/>
          <br/>
            <i>Enum</i>: whitespace, uaxUrlEmail, standard, regexSplit, regexCaptureGroup, nGram, keyword, edgeGram<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>group</b></td>
        <td>integer</td>
        <td>
          Index of the character group within the matching expression to extract into tokens. Use `0` to extract all character groups.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxGram</b></td>
        <td>integer</td>
        <td>
          Characters to include in the longest token that Atlas Search creates.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxTokenLength</b></td>
        <td>integer</td>
        <td>
          Maximum number of characters in a single token. Tokens greater than this length are split at this length into multiple tokens.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>minGram</b></td>
        <td>integer</td>
        <td>
          Characters to include in the shortest token that Atlas Search creates.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>pattern</b></td>
        <td>string</td>
        <td>
          Regular expression to match against.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasSearchIndexConfig.status
<sup><sup>[↩ Parent](#atlassearchindexconfig)</sup></sup>



AtlasSearchIndexConfigStatus defines the observed state of AtlasSearchIndexConfig.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlassearchindexconfigstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasSearchIndexConfig.status.conditions[index]
<sup><sup>[↩ Parent](#atlassearchindexconfigstatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasStreamConnection
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasStreamConnection is the Schema for the atlasstreamconnections API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasStreamConnection</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasstreamconnectionspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasStreamConnectionSpec defines the desired state of AtlasStreamConnection.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasstreamconnectionstatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasStreamConnectionStatus defines the observed state of AtlasStreamConnection.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamConnection.spec
<sup><sup>[↩ Parent](#atlasstreamconnection)</sup></sup>



AtlasStreamConnectionSpec defines the desired state of AtlasStreamConnection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that uniquely identifies the stream connection.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Type of the connection. Can be either Cluster or Kafka.<br/>
          <br/>
            <i>Enum</i>: Kafka, Cluster, Sample<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasstreamconnectionspecclusterconfig">clusterConfig</a></b></td>
        <td>object</td>
        <td>
          The configuration to be used to connect to an Atlas Cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasstreamconnectionspeckafkaconfig">kafkaConfig</a></b></td>
        <td>object</td>
        <td>
          The configuration to be used to connect to a Kafka Cluster.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamConnection.spec.clusterConfig
<sup><sup>[↩ Parent](#atlasstreamconnectionspec)</sup></sup>



The configuration to be used to connect to an Atlas Cluster.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the cluster configured for this connection.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasstreamconnectionspecclusterconfigrole">role</a></b></td>
        <td>object</td>
        <td>
          The name of a built-in or Custom DB Role to connect to an Atlas Cluster.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasStreamConnection.spec.clusterConfig.role
<sup><sup>[↩ Parent](#atlasstreamconnectionspecclusterconfig)</sup></sup>



The name of a built-in or Custom DB Role to connect to an Atlas Cluster.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The name of the role to use. Can be a built-in role or a custom role.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Type of the DB role. Can be either BUILT_IN or CUSTOM.<br/>
          <br/>
            <i>Enum</i>: BUILT_IN, CUSTOM<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasStreamConnection.spec.kafkaConfig
<sup><sup>[↩ Parent](#atlasstreamconnectionspec)</sup></sup>



The configuration to be used to connect to a Kafka Cluster.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasstreamconnectionspeckafkaconfigauthentication">authentication</a></b></td>
        <td>object</td>
        <td>
          User credentials required to connect to a Kafka Cluster. Includes the authentication type, as well as the parameters for that authentication mode.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>bootstrapServers</b></td>
        <td>string</td>
        <td>
          Comma separated list of server addresses<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasstreamconnectionspeckafkaconfigsecurity">security</a></b></td>
        <td>object</td>
        <td>
          Properties for the secure transport connection to Kafka. For SSL, this can include the trusted certificate to use.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>config</b></td>
        <td>map[string]string</td>
        <td>
          A map of Kafka key-value pairs for optional configuration. This is a flat object, and keys can have '.' characters.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamConnection.spec.kafkaConfig.authentication
<sup><sup>[↩ Parent](#atlasstreamconnectionspeckafkaconfig)</sup></sup>



User credentials required to connect to a Kafka Cluster. Includes the authentication type, as well as the parameters for that authentication mode.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasstreamconnectionspeckafkaconfigauthenticationcredentials">credentials</a></b></td>
        <td>object</td>
        <td>
          Reference to the secret containing th Username and Password of the account to connect to the Kafka cluster.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mechanism</b></td>
        <td>enum</td>
        <td>
          Style of authentication. Can be one of PLAIN, SCRAM-256, or SCRAM-512.<br/>
          <br/>
            <i>Enum</i>: PLAIN, SCRAM-256, SCRAM-512<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasStreamConnection.spec.kafkaConfig.authentication.credentials
<sup><sup>[↩ Parent](#atlasstreamconnectionspeckafkaconfigauthentication)</sup></sup>



Reference to the secret containing th Username and Password of the account to connect to the Kafka cluster.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamConnection.spec.kafkaConfig.security
<sup><sup>[↩ Parent](#atlasstreamconnectionspeckafkaconfig)</sup></sup>



Properties for the secure transport connection to Kafka. For SSL, this can include the trusted certificate to use.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>protocol</b></td>
        <td>enum</td>
        <td>
          Describes the transport type. Can be either PLAINTEXT or SSL.<br/>
          <br/>
            <i>Enum</i>: PLAINTEXT, SSL<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasstreamconnectionspeckafkaconfigsecuritycertificate">certificate</a></b></td>
        <td>object</td>
        <td>
          A trusted, public x509 certificate for connecting to Kafka over SSL.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamConnection.spec.kafkaConfig.security.certificate
<sup><sup>[↩ Parent](#atlasstreamconnectionspeckafkaconfigsecurity)</sup></sup>



A trusted, public x509 certificate for connecting to Kafka over SSL.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamConnection.status
<sup><sup>[↩ Parent](#atlasstreamconnection)</sup></sup>



AtlasStreamConnectionStatus defines the observed state of AtlasStreamConnection.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasstreamconnectionstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasstreamconnectionstatusinstancesindex">instances</a></b></td>
        <td>[]object</td>
        <td>
          List of instances using the connection configuration<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamConnection.status.conditions[index]
<sup><sup>[↩ Parent](#atlasstreamconnectionstatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamConnection.status.instances[index]
<sup><sup>[↩ Parent](#atlasstreamconnectionstatus)</sup></sup>



ResourceRefNamespaced is a reference to a Kubernetes Resource that allows to configure the namespace

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasStreamInstance
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasStreamInstance is the Schema for the atlasstreaminstances API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasStreamInstance</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasstreaminstancespec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasStreamInstanceSpec defines the desired state of AtlasStreamInstance.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasstreaminstancestatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasStreamInstanceStatus defines the observed state of AtlasStreamInstance.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamInstance.spec
<sup><sup>[↩ Parent](#atlasstreaminstance)</sup></sup>



AtlasStreamInstanceSpec defines the desired state of AtlasStreamInstance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasstreaminstancespecclusterconfig">clusterConfig</a></b></td>
        <td>object</td>
        <td>
          The configuration to be used to connect to an Atlas Cluster.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that identifies the stream connection.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasstreaminstancespecprojectref">projectRef</a></b></td>
        <td>object</td>
        <td>
          Project which the instance belongs to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasstreaminstancespecconnectionregistryindex">connectionRegistry</a></b></td>
        <td>[]object</td>
        <td>
          List of connections of the stream instance for the specified project.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamInstance.spec.clusterConfig
<sup><sup>[↩ Parent](#atlasstreaminstancespec)</sup></sup>



The configuration to be used to connect to an Atlas Cluster.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>provider</b></td>
        <td>enum</td>
        <td>
          Name of the cluster configured for this connection.<br/>
          <br/>
            <i>Enum</i>: AWS, GCP, AZURE, TENANT, SERVERLESS<br/>
            <i>Default</i>: AWS<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          Name of the cloud provider region hosting Atlas Stream Processing.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>tier</b></td>
        <td>enum</td>
        <td>
          Selected tier for the Stream Instance. Configures Memory / VCPU allowances.<br/>
          <br/>
            <i>Enum</i>: SP10, SP30, SP50<br/>
            <i>Default</i>: SP10<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasStreamInstance.spec.projectRef
<sup><sup>[↩ Parent](#atlasstreaminstancespec)</sup></sup>



Project which the instance belongs to.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamInstance.spec.connectionRegistry[index]
<sup><sup>[↩ Parent](#atlasstreaminstancespec)</sup></sup>



ResourceRefNamespaced is a reference to a Kubernetes Resource that allows to configure the namespace

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamInstance.status
<sup><sup>[↩ Parent](#atlasstreaminstance)</sup></sup>



AtlasStreamInstanceStatus defines the observed state of AtlasStreamInstance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasstreaminstancestatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasstreaminstancestatusconnectionsindex">connections</a></b></td>
        <td>[]object</td>
        <td>
          List of connections configured in the stream instance.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostnames</b></td>
        <td>[]string</td>
        <td>
          List that contains the hostnames assigned to the stream instance.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Unique 24-hexadecimal character string that identifies the instance<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamInstance.status.conditions[index]
<sup><sup>[↩ Parent](#atlasstreaminstancestatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamInstance.status.connections[index]
<sup><sup>[↩ Parent](#atlasstreaminstancestatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Human-readable label that uniquely identifies the stream connection<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasstreaminstancestatusconnectionsindexresourceref">resourceRef</a></b></td>
        <td>object</td>
        <td>
          Reference for the resource that contains connection configuration<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasStreamInstance.status.connections[index].resourceRef
<sup><sup>[↩ Parent](#atlasstreaminstancestatusconnectionsindex)</sup></sup>



Reference for the resource that contains connection configuration

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## AtlasTeam
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasTeam is the Schema for the Atlas Teams API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasTeam</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasteamspec">spec</a></b></td>
        <td>object</td>
        <td>
          TeamSpec defines the desired state of a Team in Atlas.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasteamstatus">status</a></b></td>
        <td>object</td>
        <td>
          TeamStatus defines the observed state of AtlasTeam.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasTeam.spec
<sup><sup>[↩ Parent](#atlasteam)</sup></sup>



TeamSpec defines the desired state of a Team in Atlas.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The name of the team you want to create.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>usernames</b></td>
        <td>[]string</td>
        <td>
          Valid email addresses of users to add to the new team.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasTeam.status
<sup><sup>[↩ Parent](#atlasteam)</sup></sup>



TeamStatus defines the observed state of AtlasTeam.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasteamstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions is the list of statuses showing the current state of the Atlas Custom Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID of the team<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration indicates the generation of the resource specification that the Atlas Operator is aware of.
The Atlas Operator updates this field to the 'metadata.generation' as soon as it starts reconciliation of the resource.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasteamstatusprojectsindex">projects</a></b></td>
        <td>[]object</td>
        <td>
          List of projects which the team is assigned<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasTeam.status.conditions[index]
<sup><sup>[↩ Parent](#atlasteamstatus)</sup></sup>



Condition describes the state of an Atlas Custom Resource at a certain point.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of the condition, one of True, False, Unknown.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of Atlas Custom Resource condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.
Represented in ISO 8601 format.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A human readable message indicating details about the transition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          The reason for the condition's last transition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasTeam.status.projects[index]
<sup><sup>[↩ Parent](#atlasteamstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Unique identifier of the project inside atlas<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name given to the project<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>

## AtlasThirdPartyIntegration
<sup><sup>[↩ Parent](#atlasmongodbcomv1 )</sup></sup>






AtlasThirdPartyIntegration is the Schema for the atlas 3rd party integrations API.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>atlas.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>AtlasThirdPartyIntegration</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspec">spec</a></b></td>
        <td>object</td>
        <td>
          AtlasThirdPartyIntegrationSpec contains the expected configuration for an integration<br/>
          <br/>
            <i>Validations</i>:<li>(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef)): must define only one project reference through externalProjectRef or projectRef</li><li>(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef): must define a local connection secret when referencing an external project</li><li>has(self.type) && self.type.size() != 0: must define a type of integration</li><li>!has(self.datadog) || (self.type == 'DATADOG' && has(self.datadog)): only DATADOG type may set datadog fields</li><li>!has(self.microsoftTeams) || (self.type == 'MICROSOFT_TEAMS' && has(self.microsoftTeams)): only MICROSOFT_TEAMS type may set microsoftTeams fields</li><li>!has(self.newRelic) || (self.type == 'NEW_RELIC' && has(self.newRelic)): only NEW_RELIC type may set newRelic fields</li><li>!has(self.opsGenie) || (self.type == 'OPS_GENIE' && has(self.opsGenie)): only OPS_GENIE type may set opsGenie fields</li><li>!has(self.prometheus) || (self.type == 'PROMETHEUS' && has(self.prometheus)): only PROMETHEUS type may set prometheus fields</li><li>!has(self.pagerDuty) || (self.type == 'PAGER_DUTY' && has(self.pagerDuty)): only PAGER_DUTY type may set pagerDuty fields</li><li>!has(self.slack) || (self.type == 'SLACK' && has(self.slack)): only SLACK type may set slack fields</li><li>!has(self.victorOps) || (self.type == 'VICTOR_OPS' && has(self.victorOps)): only VICTOR_OPS type may set victorOps fields</li><li>!has(self.webhook) || (self.type == 'WEBHOOK' && has(self.webhook)): only WEBHOOK type may set webhook fields</li>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationstatus">status</a></b></td>
        <td>object</td>
        <td>
          AtlasThirdPartyIntegrationStatus holds the status of an integration<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec
<sup><sup>[↩ Parent](#atlasthirdpartyintegration)</sup></sup>



AtlasThirdPartyIntegrationSpec contains the expected configuration for an integration

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Type of the integration.<br/>
          <br/>
            <i>Enum</i>: DATADOG, MICROSOFT_TEAMS, NEW_RELIC, OPS_GENIE, PAGER_DUTY, PROMETHEUS, SLACK, VICTOR_OPS, WEBHOOK<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecconnectionsecret">connectionSecret</a></b></td>
        <td>object</td>
        <td>
          Name of the secret containing Atlas API private and public keys.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecdatadog">datadog</a></b></td>
        <td>object</td>
        <td>
          Datadog contains the config fields for Datadog's Integration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecexternalprojectref">externalProjectRef</a></b></td>
        <td>object</td>
        <td>
          externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecmicrosoftteams">microsoftTeams</a></b></td>
        <td>object</td>
        <td>
          MicrosoftTeams contains the config fields for Microsoft Teams's Integration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecnewrelic">newRelic</a></b></td>
        <td>object</td>
        <td>
          NewRelic contains the config fields for New Relic's Integration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecopsgenie">opsGenie</a></b></td>
        <td>object</td>
        <td>
          OpsGenie contains the config fields for Ops Genie's Integration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecpagerduty">pagerDuty</a></b></td>
        <td>object</td>
        <td>
          PagerDuty contains the config fields for PagerDuty's Integration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecprojectref">projectRef</a></b></td>
        <td>object</td>
        <td>
          projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecprometheus">prometheus</a></b></td>
        <td>object</td>
        <td>
          Prometheus contains the config fields for Prometheus's Integration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecslack">slack</a></b></td>
        <td>object</td>
        <td>
          Slack contains the config fields for Slack's Integration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecvictorops">victorOps</a></b></td>
        <td>object</td>
        <td>
          VictorOps contains the config fields for VictorOps's Integration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecwebhook">webhook</a></b></td>
        <td>object</td>
        <td>
          Webhook contains the config fields for Webhook's Integration.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.connectionSecret
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



Name of the secret containing Atlas API private and public keys.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.datadog
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



Datadog contains the config fields for Datadog's Integration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecdatadogapikeysecretref">apiKeySecretRef</a></b></td>
        <td>object</td>
        <td>
          APIKeySecretRef holds the name of a secret containing the Datadog API key.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          Region is the Datadog region<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>sendCollectionLatencyMetrics</b></td>
        <td>enum</td>
        <td>
          SendCollectionLatencyMetrics toggles sending collection latency metrics.<br/>
          <br/>
            <i>Enum</i>: enabled, disabled<br/>
            <i>Default</i>: disabled<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sendDatabaseMetrics</b></td>
        <td>enum</td>
        <td>
          SendDatabaseMetrics toggles sending database metrics,
including database and collection names<br/>
          <br/>
            <i>Enum</i>: enabled, disabled<br/>
            <i>Default</i>: disabled<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.datadog.apiKeySecretRef
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspecdatadog)</sup></sup>



APIKeySecretRef holds the name of a secret containing the Datadog API key.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.externalProjectRef
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



externalProjectRef holds the parent Atlas project ID.
Mutually exclusive with the "projectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID is the Atlas project ID.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.microsoftTeams
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



MicrosoftTeams contains the config fields for Microsoft Teams's Integration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecmicrosoftteamsurlsecretref">urlSecretRef</a></b></td>
        <td>object</td>
        <td>
          URLSecretRef holds the name of a secret containing the Microsoft Teams secret URL.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.microsoftTeams.urlSecretRef
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspecmicrosoftteams)</sup></sup>



URLSecretRef holds the name of a secret containing the Microsoft Teams secret URL.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.newRelic
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



NewRelic contains the config fields for New Relic's Integration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecnewreliccredentialssecretref">credentialsSecretRef</a></b></td>
        <td>object</td>
        <td>
          CredentialsSecretRef holds the name of a secret containing new relic's credentials:
account id, license key, read and write tokens.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.newRelic.credentialsSecretRef
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspecnewrelic)</sup></sup>



CredentialsSecretRef holds the name of a secret containing new relic's credentials:
account id, license key, read and write tokens.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.opsGenie
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



OpsGenie contains the config fields for Ops Genie's Integration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecopsgenieapikeysecretref">apiKeySecretRef</a></b></td>
        <td>object</td>
        <td>
          APIKeySecretRef holds the name of a secret containing Ops Genie's API key.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          Region is the Ops Genie region.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.opsGenie.apiKeySecretRef
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspecopsgenie)</sup></sup>



APIKeySecretRef holds the name of a secret containing Ops Genie's API key.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.pagerDuty
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



PagerDuty contains the config fields for PagerDuty's Integration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>region</b></td>
        <td>string</td>
        <td>
          Region is the Pager Duty region.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecpagerdutyservicekeysecretref">serviceKeySecretRef</a></b></td>
        <td>object</td>
        <td>
          ServiceKeySecretRef holds the name of a secret containing Pager Duty service key.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.pagerDuty.serviceKeySecretRef
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspecpagerduty)</sup></sup>



ServiceKeySecretRef holds the name of a secret containing Pager Duty service key.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.projectRef
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



projectRef is a reference to the parent AtlasProject resource.
Mutually exclusive with the "externalProjectRef" field.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the Kubernetes Resource<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the Kubernetes Resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.prometheus
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



Prometheus contains the config fields for Prometheus's Integration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>string</td>
        <td>
          Enabled is true when Prometheus integration is enabled.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecprometheusprometheuscredentialssecretref">prometheusCredentialsSecretRef</a></b></td>
        <td>object</td>
        <td>
          PrometheusCredentialsSecretRef holds the name of a secret containing the Prometheus.
username & password<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>serviceDiscovery</b></td>
        <td>enum</td>
        <td>
          ServiceDiscovery to be used by Prometheus.<br/>
          <br/>
            <i>Enum</i>: file, http<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.prometheus.prometheusCredentialsSecretRef
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspecprometheus)</sup></sup>



PrometheusCredentialsSecretRef holds the name of a secret containing the Prometheus.
username & password

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.slack
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



Slack contains the config fields for Slack's Integration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecslackapitokensecretref">apiTokenSecretRef</a></b></td>
        <td>object</td>
        <td>
          APITokenSecretRef holds the name of a secret containing the Slack API token.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>channelName</b></td>
        <td>string</td>
        <td>
          ChannelName to be used by Prometheus.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>teamName</b></td>
        <td>string</td>
        <td>
          TeamName flags whether Prometheus integration is enabled.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.slack.apiTokenSecretRef
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspecslack)</sup></sup>



APITokenSecretRef holds the name of a secret containing the Slack API token.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.victorOps
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



VictorOps contains the config fields for VictorOps's Integration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecvictoropsapikeysecretref">apiKeySecretRef</a></b></td>
        <td>object</td>
        <td>
          APIKeySecretRef is the name of a secret containing Victor Ops API key.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>routingKey</b></td>
        <td>string</td>
        <td>
          RoutingKey holds VictorOps routing key.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.victorOps.apiKeySecretRef
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspecvictorops)</sup></sup>



APIKeySecretRef is the name of a secret containing Victor Ops API key.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.webhook
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspec)</sup></sup>



Webhook contains the config fields for Webhook's Integration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasthirdpartyintegrationspecwebhookurlsecretref">urlSecretRef</a></b></td>
        <td>object</td>
        <td>
          URLSecretRef holds the name of a secret containing Webhook URL and secret.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.spec.webhook.urlSecretRef
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationspecwebhook)</sup></sup>



URLSecretRef holds the name of a secret containing Webhook URL and secret.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource being referred to
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.status
<sup><sup>[↩ Parent](#atlasthirdpartyintegration)</sup></sup>



AtlasThirdPartyIntegrationStatus holds the status of an integration

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#atlasthirdpartyintegrationstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions holding the status details<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID of the third party integration resource in Atlas<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### AtlasThirdPartyIntegration.status.conditions[index]
<sup><sup>[↩ Parent](#atlasthirdpartyintegrationstatus)</sup></sup>



Condition contains details for one aspect of the current state of this API Resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          lastTransitionTime is the last time the condition transitioned from one status to another.
This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          message is a human readable message indicating details about the transition.
This may be an empty string.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          reason contains a programmatic identifier indicating the reason for the condition's last transition.
Producers of specific condition types may define expected values and meanings for this field,
and whether the values are considered a guaranteed API.
The value should be a CamelCase string.
This field may not be empty.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>enum</td>
        <td>
          status of the condition, one of True, False, Unknown.<br/>
          <br/>
            <i>Enum</i>: True, False, Unknown<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type of condition in CamelCase or in foo.example.com/CamelCase.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          observedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Minimum</i>: 0<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
