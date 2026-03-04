# API Reference

Packages:

- [atlas.generated.mongodb.com/v1](#atlasgeneratedmongodbcomv1)

# atlas.generated.mongodb.com/v1

Resource Types:

- [Group](#group)




## Group
<sup><sup>[↩ Parent](#atlasgeneratedmongodbcomv1 )</sup></sup>






A group, managed by the MongoDB Kubernetes Atlas Operator.

<table>
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
      <td>atlas.generated.mongodb.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>Group</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#groupspec">spec</a></b></td>
        <td>object</td>
        <td>
          Specification of the group supporting the following versions:

- v20250312

At most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#groupstatus">status</a></b></td>
        <td>object</td>
        <td>
          Most recently observed read-only status of the group for the specified resource version. This data may not be up to date and is populated by the system. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Group.spec
<sup><sup>[↩ Parent](#group)</sup></sup>



Specification of the group supporting the following versions:

- v20250312

At most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#groupspecconnectionsecretref">connectionSecretRef</a></b></td>
        <td>object</td>
        <td>
          SENSITIVE FIELD

Reference to a secret containing the credentials to setup the connection to Atlas.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#groupspecv20250312">v20250312</a></b></td>
        <td>object</td>
        <td>
          The spec of the group resource for version v20250312.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Group.spec.connectionSecretRef
<sup><sup>[↩ Parent](#groupspec)</sup></sup>



SENSITIVE FIELD

Reference to a secret containing the credentials to setup the connection to Atlas.

<table>
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
          Name of the secret containing the Atlas credentials.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Group.spec.v20250312
<sup><sup>[↩ Parent](#groupspec)</sup></sup>



The spec of the group resource for version v20250312.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#groupspecv20250312entry">entry</a></b></td>
        <td>object</td>
        <td>
          The entry fields of the group resource spec. These fields can be set for creating and updating groups.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>projectOwnerId</b></td>
        <td>string</td>
        <td>
          Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user to whom to grant the Project Owner role on the specified project. If you set this parameter, it overrides the default value of the oldest Organization Owner.<br/>
          <br/>
            <i>Validations</i>:<li>self == oldSelf: projectOwnerId cannot be modified after creation</li>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Group.spec.v20250312.entry
<sup><sup>[↩ Parent](#groupspecv20250312)</sup></sup>



The entry fields of the group resource spec. These fields can be set for creating and updating groups.

<table>
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
          Human-readable label that identifies the project included in the MongoDB Cloud organization.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>orgId</b></td>
        <td>string</td>
        <td>
          Unique 24-hexadecimal digit string that identifies the MongoDB Cloud organization to which the project belongs.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>regionUsageRestrictions</b></td>
        <td>string</td>
        <td>
          Applies to Atlas for Government only.

In Commercial Atlas, this field will be rejected in requests and missing in responses.

This field sets restrictions on available regions in the project.

`COMMERCIAL_FEDRAMP_REGIONS_ONLY`: Only allows deployments in FedRAMP Moderate regions.

`GOV_REGIONS_ONLY`: Only allows deployments in GovCloud regions.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#groupspecv20250312entrytagsindex">tags</a></b></td>
        <td>[]object</td>
        <td>
          List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the project.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>withDefaultAlertsSettings</b></td>
        <td>boolean</td>
        <td>
          Flag that indicates whether to create the project with default alert settings. This setting cannot be updated after project creation.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Group.spec.v20250312.entry.tags[index]
<sup><sup>[↩ Parent](#groupspecv20250312entry)</sup></sup>



Key-value pair that tags and categorizes a MongoDB Cloud organization, project, or cluster. For example, `environment : production`.

<table>
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
          Constant that defines the set of the tag. For example, `environment` in the `environment : production` tag.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Variable that belongs to the set of the tag. For example, `production` in the `environment : production` tag.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Group.status
<sup><sup>[↩ Parent](#group)</sup></sup>



Most recently observed read-only status of the group for the specified resource version. This data may not be up to date and is populated by the system. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#groupstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Represents the latest available observations of a resource's current state.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#groupstatusv20250312">v20250312</a></b></td>
        <td>object</td>
        <td>
          The last observed Atlas state of the group resource for version v20250312.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Group.status.conditions[index]
<sup><sup>[↩ Parent](#groupstatus)</sup></sup>





<table>
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
          Type of condition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          Last time the condition transitioned from one status to another.<br/>
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
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          observedGeneration represents the .metadata.generation that the condition was set based upon.<br/>
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


### Group.status.v20250312
<sup><sup>[↩ Parent](#groupstatus)</sup></sup>



The last observed Atlas state of the group resource for version v20250312.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>clusterCount</b></td>
        <td>integer</td>
        <td>
          Quantity of MongoDB Cloud clusters deployed in this project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>created</b></td>
        <td>string</td>
        <td>
          Date and time when MongoDB Cloud created this project. This parameter expresses its value in the ISO 8601 timestamp format in UTC.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Unique 24-hexadecimal digit string that identifies the MongoDB Cloud project.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
