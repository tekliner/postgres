package eventgrid

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest/date"
)

// JobState enumerates the values for job state.
type JobState string

const (
	// Canceled The job was canceled. This is a final state for the job.
	Canceled JobState = "Canceled"
	// Canceling The job is in the process of being canceled. This is a transient state for the job.
	Canceling JobState = "Canceling"
	// Error The job has encountered an error. This is a final state for the job.
	Error JobState = "Error"
	// Finished The job is finished. This is a final state for the job.
	Finished JobState = "Finished"
	// Processing The job is processing. This is a transient state for the job.
	Processing JobState = "Processing"
	// Queued The job is in a queued state, waiting for resources to become available. This is a transient
	// state.
	Queued JobState = "Queued"
	// Scheduled The job is being scheduled to run on an available resource. This is a transient state, between
	// queued and processing states.
	Scheduled JobState = "Scheduled"
)

// PossibleJobStateValues returns an array of possible values for the JobState const type.
func PossibleJobStateValues() []JobState {
	return []JobState{Canceled, Canceling, Error, Finished, Processing, Queued, Scheduled}
}

// ContainerRegistryEventActor the agent that initiated the event. For most situations, this could be from the
// authorization context of the request.
type ContainerRegistryEventActor struct {
	// Name - The subject or username associated with the request context that generated the event.
	Name *string `json:"name,omitempty"`
}

// ContainerRegistryEventData the content of the event request message.
type ContainerRegistryEventData struct {
	// ID - The event ID.
	ID *string `json:"id,omitempty"`
	// Timestamp - The time at which the event occurred.
	Timestamp *date.Time `json:"timestamp,omitempty"`
	// Action - The action that encompasses the provided event.
	Action *string `json:"action,omitempty"`
	// Target - The target of the event.
	Target *ContainerRegistryEventTarget `json:"target,omitempty"`
	// Request - The request that generated the event.
	Request *ContainerRegistryEventRequest `json:"request,omitempty"`
	// Actor - The agent that initiated the event. For most situations, this could be from the authorization context of the request.
	Actor *ContainerRegistryEventActor `json:"actor,omitempty"`
	// Source - The registry node that generated the event. Put differently, while the actor initiates the event, the source generates it.
	Source *ContainerRegistryEventSource `json:"source,omitempty"`
}

// ContainerRegistryEventRequest the request that generated the event.
type ContainerRegistryEventRequest struct {
	// ID - The ID of the request that initiated the event.
	ID *string `json:"id,omitempty"`
	// Addr - The IP or hostname and possibly port of the client connection that initiated the event. This is the RemoteAddr from the standard http request.
	Addr *string `json:"addr,omitempty"`
	// Host - The externally accessible hostname of the registry instance, as specified by the http host header on incoming requests.
	Host *string `json:"host,omitempty"`
	// Method - The request method that generated the event.
	Method *string `json:"method,omitempty"`
	// Useragent - The user agent header of the request.
	Useragent *string `json:"useragent,omitempty"`
}

// ContainerRegistryEventSource the registry node that generated the event. Put differently, while the actor
// initiates the event, the source generates it.
type ContainerRegistryEventSource struct {
	// Addr - The IP or hostname and the port of the registry node that generated the event. Generally, this will be resolved by os.Hostname() along with the running port.
	Addr *string `json:"addr,omitempty"`
	// InstanceID - The running instance of an application. Changes after each restart.
	InstanceID *string `json:"instanceID,omitempty"`
}

// ContainerRegistryEventTarget the target of the event.
type ContainerRegistryEventTarget struct {
	// MediaType - The MIME type of the referenced object.
	MediaType *string `json:"mediaType,omitempty"`
	// Size - The number of bytes of the content. Same as Length field.
	Size *int64 `json:"size,omitempty"`
	// Digest - The digest of the content, as defined by the Registry V2 HTTP API Specification.
	Digest *string `json:"digest,omitempty"`
	// Length - The number of bytes of the content. Same as Size field.
	Length *int64 `json:"length,omitempty"`
	// Repository - The repository name.
	Repository *string `json:"repository,omitempty"`
	// URL - The direct URL to the content.
	URL *string `json:"url,omitempty"`
	// Tag - The tag name.
	Tag *string `json:"tag,omitempty"`
}

// ContainerRegistryImageDeletedEventData schema of the Data property of an EventGridEvent for a
// Microsoft.ContainerRegistry.ImageDeleted event.
type ContainerRegistryImageDeletedEventData struct {
	// ID - The event ID.
	ID *string `json:"id,omitempty"`
	// Timestamp - The time at which the event occurred.
	Timestamp *date.Time `json:"timestamp,omitempty"`
	// Action - The action that encompasses the provided event.
	Action *string `json:"action,omitempty"`
	// Target - The target of the event.
	Target *ContainerRegistryEventTarget `json:"target,omitempty"`
	// Request - The request that generated the event.
	Request *ContainerRegistryEventRequest `json:"request,omitempty"`
	// Actor - The agent that initiated the event. For most situations, this could be from the authorization context of the request.
	Actor *ContainerRegistryEventActor `json:"actor,omitempty"`
	// Source - The registry node that generated the event. Put differently, while the actor initiates the event, the source generates it.
	Source *ContainerRegistryEventSource `json:"source,omitempty"`
}

// ContainerRegistryImagePushedEventData schema of the Data property of an EventGridEvent for a
// Microsoft.ContainerRegistry.ImagePushed event.
type ContainerRegistryImagePushedEventData struct {
	// ID - The event ID.
	ID *string `json:"id,omitempty"`
	// Timestamp - The time at which the event occurred.
	Timestamp *date.Time `json:"timestamp,omitempty"`
	// Action - The action that encompasses the provided event.
	Action *string `json:"action,omitempty"`
	// Target - The target of the event.
	Target *ContainerRegistryEventTarget `json:"target,omitempty"`
	// Request - The request that generated the event.
	Request *ContainerRegistryEventRequest `json:"request,omitempty"`
	// Actor - The agent that initiated the event. For most situations, this could be from the authorization context of the request.
	Actor *ContainerRegistryEventActor `json:"actor,omitempty"`
	// Source - The registry node that generated the event. Put differently, while the actor initiates the event, the source generates it.
	Source *ContainerRegistryEventSource `json:"source,omitempty"`
}

// DeviceConnectionStateEventInfo information about the device connection state event.
type DeviceConnectionStateEventInfo struct {
	// SequenceNumber - Sequence number is string representation of a hexadecimal number. string compare can be used to identify the larger number because both in ASCII and HEX numbers come after alphabets. If you are converting the string to hex, then the number is a 256 bit number.
	SequenceNumber *string `json:"sequenceNumber,omitempty"`
}

// DeviceConnectionStateEventProperties schema of the Data property of an EventGridEvent for a device connection
// state event (DeviceConnected, DeviceDisconnected).
type DeviceConnectionStateEventProperties struct {
	// DeviceID - The unique identifier of the device. This case-sensitive string can be up to 128 characters long, and supports ASCII 7-bit alphanumeric characters plus the following special characters: - : . + % _ &#35; * ? ! ( ) , = @ ; $ '.
	DeviceID *string `json:"deviceId,omitempty"`
	// ModuleID - The unique identifier of the module. This case-sensitive string can be up to 128 characters long, and supports ASCII 7-bit alphanumeric characters plus the following special characters: - : . + % _ &#35; * ? ! ( ) , = @ ; $ '.
	ModuleID *string `json:"moduleId,omitempty"`
	// HubName - Name of the IoT Hub where the device was created or deleted.
	HubName *string `json:"hubName,omitempty"`
	// DeviceConnectionStateEventInfo - Information about the device connection state event.
	DeviceConnectionStateEventInfo *DeviceConnectionStateEventInfo `json:"deviceConnectionStateEventInfo,omitempty"`
}

// DeviceLifeCycleEventProperties schema of the Data property of an EventGridEvent for a device life cycle event
// (DeviceCreated, DeviceDeleted).
type DeviceLifeCycleEventProperties struct {
	// DeviceID - The unique identifier of the device. This case-sensitive string can be up to 128 characters long, and supports ASCII 7-bit alphanumeric characters plus the following special characters: - : . + % _ &#35; * ? ! ( ) , = @ ; $ '.
	DeviceID *string `json:"deviceId,omitempty"`
	// HubName - Name of the IoT Hub where the device was created or deleted.
	HubName *string `json:"hubName,omitempty"`
	// Twin - Information about the device twin, which is the cloud representation of application device metadata.
	Twin *DeviceTwinInfo `json:"twin,omitempty"`
}

// DeviceTwinInfo information about the device twin, which is the cloud representation of application device
// metadata.
type DeviceTwinInfo struct {
	// AuthenticationType - Authentication type used for this device: either SAS, SelfSigned, or CertificateAuthority.
	AuthenticationType *string `json:"authenticationType,omitempty"`
	// CloudToDeviceMessageCount - Count of cloud to device messages sent to this device.
	CloudToDeviceMessageCount *float64 `json:"cloudToDeviceMessageCount,omitempty"`
	// ConnectionState - Whether the device is connected or disconnected.
	ConnectionState *string `json:"connectionState,omitempty"`
	// DeviceID - The unique identifier of the device twin.
	DeviceID *string `json:"deviceId,omitempty"`
	// Etag - A piece of information that describes the content of the device twin. Each etag is guaranteed to be unique per device twin.
	Etag *string `json:"etag,omitempty"`
	// LastActivityTime - The ISO8601 timestamp of the last activity.
	LastActivityTime *string `json:"lastActivityTime,omitempty"`
	// Properties - Properties JSON element.
	Properties *DeviceTwinInfoProperties `json:"properties,omitempty"`
	// Status - Whether the device twin is enabled or disabled.
	Status *string `json:"status,omitempty"`
	// StatusUpdateTime - The ISO8601 timestamp of the last device twin status update.
	StatusUpdateTime *string `json:"statusUpdateTime,omitempty"`
	// Version - An integer that is incremented by one each time the device twin is updated.
	Version *float64 `json:"version,omitempty"`
	// X509Thumbprint - The thumbprint is a unique value for the x509 certificate, commonly used to find a particular certificate in a certificate store. The thumbprint is dynamically generated using the SHA1 algorithm, and does not physically exist in the certificate.
	X509Thumbprint *DeviceTwinInfoX509Thumbprint `json:"x509Thumbprint,omitempty"`
}

// DeviceTwinInfoProperties properties JSON element.
type DeviceTwinInfoProperties struct {
	// Desired - A portion of the properties that can be written only by the application back-end, and read by the device.
	Desired *DeviceTwinProperties `json:"desired,omitempty"`
	// Reported - A portion of the properties that can be written only by the device, and read by the application back-end.
	Reported *DeviceTwinProperties `json:"reported,omitempty"`
}

// DeviceTwinInfoX509Thumbprint the thumbprint is a unique value for the x509 certificate, commonly used to find a
// particular certificate in a certificate store. The thumbprint is dynamically generated using the SHA1 algorithm,
// and does not physically exist in the certificate.
type DeviceTwinInfoX509Thumbprint struct {
	// PrimaryThumbprint - Primary thumbprint for the x509 certificate.
	PrimaryThumbprint *string `json:"primaryThumbprint,omitempty"`
	// SecondaryThumbprint - Secondary thumbprint for the x509 certificate.
	SecondaryThumbprint *string `json:"secondaryThumbprint,omitempty"`
}

// DeviceTwinMetadata metadata information for the properties JSON document.
type DeviceTwinMetadata struct {
	// LastUpdated - The ISO8601 timestamp of the last time the properties were updated.
	LastUpdated *string `json:"lastUpdated,omitempty"`
}

// DeviceTwinProperties a portion of the properties that can be written only by the application back-end, and read
// by the device.
type DeviceTwinProperties struct {
	// Metadata - Metadata information for the properties JSON document.
	Metadata *DeviceTwinMetadata `json:"metadata,omitempty"`
	// Version - Version of device twin properties.
	Version *float64 `json:"version,omitempty"`
}

// Event properties of an event published to an Event Grid topic.
type Event struct {
	// ID - An unique identifier for the event.
	ID *string `json:"id,omitempty"`
	// Topic - The resource path of the event source.
	Topic *string `json:"topic,omitempty"`
	// Subject - A resource path relative to the topic path.
	Subject *string `json:"subject,omitempty"`
	// Data - Event data specific to the event type.
	Data interface{} `json:"data,omitempty"`
	// EventType - The type of the event that occurred.
	EventType *string `json:"eventType,omitempty"`
	// EventTime - The time (in UTC) the event was generated.
	EventTime *date.Time `json:"eventTime,omitempty"`
	// MetadataVersion - The schema version of the event metadata.
	MetadataVersion *string `json:"metadataVersion,omitempty"`
	// DataVersion - The schema version of the data object.
	DataVersion *string `json:"dataVersion,omitempty"`
}

// EventHubCaptureFileCreatedEventData schema of the Data property of an EventGridEvent for an
// Microsoft.EventHub.CaptureFileCreated event.
type EventHubCaptureFileCreatedEventData struct {
	// Fileurl - The path to the capture file.
	Fileurl *string `json:"fileurl,omitempty"`
	// FileType - The file type of the capture file.
	FileType *string `json:"fileType,omitempty"`
	// PartitionID - The shard ID.
	PartitionID *string `json:"partitionId,omitempty"`
	// SizeInBytes - The file size.
	SizeInBytes *int32 `json:"sizeInBytes,omitempty"`
	// EventCount - The number of events in the file.
	EventCount *int32 `json:"eventCount,omitempty"`
	// FirstSequenceNumber - The smallest sequence number from the queue.
	FirstSequenceNumber *int32 `json:"firstSequenceNumber,omitempty"`
	// LastSequenceNumber - The last sequence number from the queue.
	LastSequenceNumber *int32 `json:"lastSequenceNumber,omitempty"`
	// FirstEnqueueTime - The first time from the queue.
	FirstEnqueueTime *date.Time `json:"firstEnqueueTime,omitempty"`
	// LastEnqueueTime - The last time from the queue.
	LastEnqueueTime *date.Time `json:"lastEnqueueTime,omitempty"`
}

// IotHubDeviceConnectedEventData event data for Microsoft.Devices.DeviceConnected event.
type IotHubDeviceConnectedEventData struct {
	// DeviceID - The unique identifier of the device. This case-sensitive string can be up to 128 characters long, and supports ASCII 7-bit alphanumeric characters plus the following special characters: - : . + % _ &#35; * ? ! ( ) , = @ ; $ '.
	DeviceID *string `json:"deviceId,omitempty"`
	// ModuleID - The unique identifier of the module. This case-sensitive string can be up to 128 characters long, and supports ASCII 7-bit alphanumeric characters plus the following special characters: - : . + % _ &#35; * ? ! ( ) , = @ ; $ '.
	ModuleID *string `json:"moduleId,omitempty"`
	// HubName - Name of the IoT Hub where the device was created or deleted.
	HubName *string `json:"hubName,omitempty"`
	// DeviceConnectionStateEventInfo - Information about the device connection state event.
	DeviceConnectionStateEventInfo *DeviceConnectionStateEventInfo `json:"deviceConnectionStateEventInfo,omitempty"`
}

// IotHubDeviceCreatedEventData event data for Microsoft.Devices.DeviceCreated event.
type IotHubDeviceCreatedEventData struct {
	// DeviceID - The unique identifier of the device. This case-sensitive string can be up to 128 characters long, and supports ASCII 7-bit alphanumeric characters plus the following special characters: - : . + % _ &#35; * ? ! ( ) , = @ ; $ '.
	DeviceID *string `json:"deviceId,omitempty"`
	// HubName - Name of the IoT Hub where the device was created or deleted.
	HubName *string `json:"hubName,omitempty"`
	// Twin - Information about the device twin, which is the cloud representation of application device metadata.
	Twin *DeviceTwinInfo `json:"twin,omitempty"`
}

// IotHubDeviceDeletedEventData event data for Microsoft.Devices.DeviceDeleted event.
type IotHubDeviceDeletedEventData struct {
	// DeviceID - The unique identifier of the device. This case-sensitive string can be up to 128 characters long, and supports ASCII 7-bit alphanumeric characters plus the following special characters: - : . + % _ &#35; * ? ! ( ) , = @ ; $ '.
	DeviceID *string `json:"deviceId,omitempty"`
	// HubName - Name of the IoT Hub where the device was created or deleted.
	HubName *string `json:"hubName,omitempty"`
	// Twin - Information about the device twin, which is the cloud representation of application device metadata.
	Twin *DeviceTwinInfo `json:"twin,omitempty"`
}

// IotHubDeviceDisconnectedEventData event data for Microsoft.Devices.DeviceDisconnected event.
type IotHubDeviceDisconnectedEventData struct {
	// DeviceID - The unique identifier of the device. This case-sensitive string can be up to 128 characters long, and supports ASCII 7-bit alphanumeric characters plus the following special characters: - : . + % _ &#35; * ? ! ( ) , = @ ; $ '.
	DeviceID *string `json:"deviceId,omitempty"`
	// ModuleID - The unique identifier of the module. This case-sensitive string can be up to 128 characters long, and supports ASCII 7-bit alphanumeric characters plus the following special characters: - : . + % _ &#35; * ? ! ( ) , = @ ; $ '.
	ModuleID *string `json:"moduleId,omitempty"`
	// HubName - Name of the IoT Hub where the device was created or deleted.
	HubName *string `json:"hubName,omitempty"`
	// DeviceConnectionStateEventInfo - Information about the device connection state event.
	DeviceConnectionStateEventInfo *DeviceConnectionStateEventInfo `json:"deviceConnectionStateEventInfo,omitempty"`
}

// MediaJobStateChangeEventData schema of the Data property of an EventGridEvent for a
// Microsoft.Media.JobStateChange event.
type MediaJobStateChangeEventData struct {
	// PreviousState - The previous state of the Job. Possible values include: 'Canceled', 'Canceling', 'Error', 'Finished', 'Processing', 'Queued', 'Scheduled'
	PreviousState JobState `json:"previousState,omitempty"`
	// State - The new state of the Job. Possible values include: 'Canceled', 'Canceling', 'Error', 'Finished', 'Processing', 'Queued', 'Scheduled'
	State JobState `json:"state,omitempty"`
}

// ResourceDeleteCancelData schema of the Data property of an EventGridEvent for an
// Microsoft.Resources.ResourceDeleteCancel event. This is raised when a resource delete operation is canceled.
type ResourceDeleteCancelData struct {
	// TenantID - The tenant ID of the resource.
	TenantID *string `json:"tenantId,omitempty"`
	// SubscriptionID - The subscription ID of the resource.
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	// ResourceGroup - The resource group of the resource.
	ResourceGroup *string `json:"resourceGroup,omitempty"`
	// ResourceProvider - The resource provider performing the operation.
	ResourceProvider *string `json:"resourceProvider,omitempty"`
	// ResourceURI - The URI of the resource in the operation.
	ResourceURI *string `json:"resourceUri,omitempty"`
	// OperationName - The operation that was performed.
	OperationName *string `json:"operationName,omitempty"`
	// Status - The status of the operation.
	Status *string `json:"status,omitempty"`
	// Authorization - The requested authorization for the operation.
	Authorization *string `json:"authorization,omitempty"`
	// Claims - The properties of the claims.
	Claims *string `json:"claims,omitempty"`
	// CorrelationID - An operation ID used for troubleshooting.
	CorrelationID *string `json:"correlationId,omitempty"`
	// HTTPRequest - The details of the operation.
	HTTPRequest *string `json:"httpRequest,omitempty"`
}

// ResourceDeleteFailureData schema of the Data property of an EventGridEvent for a
// Microsoft.Resources.ResourceDeleteFailure event. This is raised when a resource delete operation fails.
type ResourceDeleteFailureData struct {
	// TenantID - The tenant ID of the resource.
	TenantID *string `json:"tenantId,omitempty"`
	// SubscriptionID - The subscription ID of the resource.
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	// ResourceGroup - The resource group of the resource.
	ResourceGroup *string `json:"resourceGroup,omitempty"`
	// ResourceProvider - The resource provider performing the operation.
	ResourceProvider *string `json:"resourceProvider,omitempty"`
	// ResourceURI - The URI of the resource in the operation.
	ResourceURI *string `json:"resourceUri,omitempty"`
	// OperationName - The operation that was performed.
	OperationName *string `json:"operationName,omitempty"`
	// Status - The status of the operation.
	Status *string `json:"status,omitempty"`
	// Authorization - The requested authorization for the operation.
	Authorization *string `json:"authorization,omitempty"`
	// Claims - The properties of the claims.
	Claims *string `json:"claims,omitempty"`
	// CorrelationID - An operation ID used for troubleshooting.
	CorrelationID *string `json:"correlationId,omitempty"`
	// HTTPRequest - The details of the operation.
	HTTPRequest *string `json:"httpRequest,omitempty"`
}

// ResourceDeleteSuccessData schema of the Data property of an EventGridEvent for a
// Microsoft.Resources.ResourceDeleteSuccess event. This is raised when a resource delete operation succeeds.
type ResourceDeleteSuccessData struct {
	// TenantID - The tenant ID of the resource.
	TenantID *string `json:"tenantId,omitempty"`
	// SubscriptionID - The subscription ID of the resource.
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	// ResourceGroup - The resource group of the resource.
	ResourceGroup *string `json:"resourceGroup,omitempty"`
	// ResourceProvider - The resource provider performing the operation.
	ResourceProvider *string `json:"resourceProvider,omitempty"`
	// ResourceURI - The URI of the resource in the operation.
	ResourceURI *string `json:"resourceUri,omitempty"`
	// OperationName - The operation that was performed.
	OperationName *string `json:"operationName,omitempty"`
	// Status - The status of the operation.
	Status *string `json:"status,omitempty"`
	// Authorization - The requested authorization for the operation.
	Authorization *string `json:"authorization,omitempty"`
	// Claims - The properties of the claims.
	Claims *string `json:"claims,omitempty"`
	// CorrelationID - An operation ID used for troubleshooting.
	CorrelationID *string `json:"correlationId,omitempty"`
	// HTTPRequest - The details of the operation.
	HTTPRequest *string `json:"httpRequest,omitempty"`
}

// ResourceWriteCancelData schema of the Data property of an EventGridEvent for a
// Microsoft.Resources.ResourceWriteCancel event. This is raised when a resource create or update operation is
// canceled.
type ResourceWriteCancelData struct {
	// TenantID - The tenant ID of the resource.
	TenantID *string `json:"tenantId,omitempty"`
	// SubscriptionID - The subscription ID of the resource.
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	// ResourceGroup - The resource group of the resource.
	ResourceGroup *string `json:"resourceGroup,omitempty"`
	// ResourceProvider - The resource provider performing the operation.
	ResourceProvider *string `json:"resourceProvider,omitempty"`
	// ResourceURI - The URI of the resource in the operation.
	ResourceURI *string `json:"resourceUri,omitempty"`
	// OperationName - The operation that was performed.
	OperationName *string `json:"operationName,omitempty"`
	// Status - The status of the operation.
	Status *string `json:"status,omitempty"`
	// Authorization - The requested authorization for the operation.
	Authorization *string `json:"authorization,omitempty"`
	// Claims - The properties of the claims.
	Claims *string `json:"claims,omitempty"`
	// CorrelationID - An operation ID used for troubleshooting.
	CorrelationID *string `json:"correlationId,omitempty"`
	// HTTPRequest - The details of the operation.
	HTTPRequest *string `json:"httpRequest,omitempty"`
}

// ResourceWriteFailureData schema of the Data property of an EventGridEvent for a
// Microsoft.Resources.ResourceWriteFailure event. This is raised when a resource create or update operation fails.
type ResourceWriteFailureData struct {
	// TenantID - The tenant ID of the resource.
	TenantID *string `json:"tenantId,omitempty"`
	// SubscriptionID - The subscription ID of the resource.
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	// ResourceGroup - The resource group of the resource.
	ResourceGroup *string `json:"resourceGroup,omitempty"`
	// ResourceProvider - The resource provider performing the operation.
	ResourceProvider *string `json:"resourceProvider,omitempty"`
	// ResourceURI - The URI of the resource in the operation.
	ResourceURI *string `json:"resourceUri,omitempty"`
	// OperationName - The operation that was performed.
	OperationName *string `json:"operationName,omitempty"`
	// Status - The status of the operation.
	Status *string `json:"status,omitempty"`
	// Authorization - The requested authorization for the operation.
	Authorization *string `json:"authorization,omitempty"`
	// Claims - The properties of the claims.
	Claims *string `json:"claims,omitempty"`
	// CorrelationID - An operation ID used for troubleshooting.
	CorrelationID *string `json:"correlationId,omitempty"`
	// HTTPRequest - The details of the operation.
	HTTPRequest *string `json:"httpRequest,omitempty"`
}

// ResourceWriteSuccessData schema of the Data property of an EventGridEvent for a
// Microsoft.Resources.ResourceWriteSuccess event. This is raised when a resource create or update operation
// succeeds.
type ResourceWriteSuccessData struct {
	// TenantID - The tenant ID of the resource.
	TenantID *string `json:"tenantId,omitempty"`
	// SubscriptionID - The subscription ID of the resource.
	SubscriptionID *string `json:"subscriptionId,omitempty"`
	// ResourceGroup - The resource group of the resource.
	ResourceGroup *string `json:"resourceGroup,omitempty"`
	// ResourceProvider - The resource provider performing the operation.
	ResourceProvider *string `json:"resourceProvider,omitempty"`
	// ResourceURI - The URI of the resource in the operation.
	ResourceURI *string `json:"resourceUri,omitempty"`
	// OperationName - The operation that was performed.
	OperationName *string `json:"operationName,omitempty"`
	// Status - The status of the operation.
	Status *string `json:"status,omitempty"`
	// Authorization - The requested authorization for the operation.
	Authorization *string `json:"authorization,omitempty"`
	// Claims - The properties of the claims.
	Claims *string `json:"claims,omitempty"`
	// CorrelationID - An operation ID used for troubleshooting.
	CorrelationID *string `json:"correlationId,omitempty"`
	// HTTPRequest - The details of the operation.
	HTTPRequest *string `json:"httpRequest,omitempty"`
}

// ServiceBusActiveMessagesAvailableWithNoListenersEventData schema of the Data property of an EventGridEvent for a
// Microsoft.ServiceBus.ActiveMessagesAvailableWithNoListeners event.
type ServiceBusActiveMessagesAvailableWithNoListenersEventData struct {
	// NamespaceName - The namespace name of the Microsoft.ServiceBus resource.
	NamespaceName *string `json:"namespaceName,omitempty"`
	// RequestURI - The endpoint of the Microsoft.ServiceBus resource.
	RequestURI *string `json:"requestUri,omitempty"`
	// EntityType - The entity type of the Microsoft.ServiceBus resource. Could be one of 'queue' or 'subscriber'.
	EntityType *string `json:"entityType,omitempty"`
	// QueueName - The name of the Microsoft.ServiceBus queue. If the entity type is of type 'subscriber', then this value will be null.
	QueueName *string `json:"queueName,omitempty"`
	// TopicName - The name of the Microsoft.ServiceBus topic. If the entity type is of type 'queue', then this value will be null.
	TopicName *string `json:"topicName,omitempty"`
	// SubscriptionName - The name of the Microsoft.ServiceBus topic's subscription. If the entity type is of type 'queue', then this value will be null.
	SubscriptionName *string `json:"subscriptionName,omitempty"`
}

// ServiceBusDeadletterMessagesAvailableWithNoListenersEventData schema of the Data property of an EventGridEvent
// for a Microsoft.ServiceBus.DeadletterMessagesAvailableWithNoListenersEvent event.
type ServiceBusDeadletterMessagesAvailableWithNoListenersEventData struct {
	// NamespaceName - The namespace name of the Microsoft.ServiceBus resource.
	NamespaceName *string `json:"namespaceName,omitempty"`
	// RequestURI - The endpoint of the Microsoft.ServiceBus resource.
	RequestURI *string `json:"requestUri,omitempty"`
	// EntityType - The entity type of the Microsoft.ServiceBus resource. Could be one of 'queue' or 'subscriber'.
	EntityType *string `json:"entityType,omitempty"`
	// QueueName - The name of the Microsoft.ServiceBus queue. If the entity type is of type 'subscriber', then this value will be null.
	QueueName *string `json:"queueName,omitempty"`
	// TopicName - The name of the Microsoft.ServiceBus topic. If the entity type is of type 'queue', then this value will be null.
	TopicName *string `json:"topicName,omitempty"`
	// SubscriptionName - The name of the Microsoft.ServiceBus topic's subscription. If the entity type is of type 'queue', then this value will be null.
	SubscriptionName *string `json:"subscriptionName,omitempty"`
}

// StorageBlobCreatedEventData schema of the Data property of an EventGridEvent for an
// Microsoft.Storage.BlobCreated event.
type StorageBlobCreatedEventData struct {
	// API - The name of the API/operation that triggered this event.
	API *string `json:"api,omitempty"`
	// ClientRequestID - A request id provided by the client of the storage API operation that triggered this event.
	ClientRequestID *string `json:"clientRequestId,omitempty"`
	// RequestID - The request id generated by the Storage service for the storage API operation that triggered this event.
	RequestID *string `json:"requestId,omitempty"`
	// ETag - The etag of the object at the time this event was triggered.
	ETag *string `json:"eTag,omitempty"`
	// ContentType - The content type of the blob. This is the same as what would be returned in the Content-Type header from the blob.
	ContentType *string `json:"contentType,omitempty"`
	// ContentLength - The size of the blob in bytes. This is the same as what would be returned in the Content-Length header from the blob.
	ContentLength *int32 `json:"contentLength,omitempty"`
	// BlobType - The type of blob.
	BlobType *string `json:"blobType,omitempty"`
	// URL - The path to the blob.
	URL *string `json:"url,omitempty"`
	// Sequencer - An opaque string value representing the logical sequence of events for any particular blob name. Users can use standard string comparison to understand the relative sequence of two events on the same blob name.
	Sequencer *string `json:"sequencer,omitempty"`
	// StorageDiagnostics - For service use only. Diagnostic data occasionally included by the Azure Storage service. This property should be ignored by event consumers.
	StorageDiagnostics interface{} `json:"storageDiagnostics,omitempty"`
}

// StorageBlobDeletedEventData schema of the Data property of an EventGridEvent for an
// Microsoft.Storage.BlobDeleted event.
type StorageBlobDeletedEventData struct {
	// API - The name of the API/operation that triggered this event.
	API *string `json:"api,omitempty"`
	// ClientRequestID - A request id provided by the client of the storage API operation that triggered this event.
	ClientRequestID *string `json:"clientRequestId,omitempty"`
	// RequestID - The request id generated by the Storage service for the storage API operation that triggered this event.
	RequestID *string `json:"requestId,omitempty"`
	// ContentType - The content type of the blob. This is the same as what would be returned in the Content-Type header from the blob.
	ContentType *string `json:"contentType,omitempty"`
	// BlobType - The type of blob.
	BlobType *string `json:"blobType,omitempty"`
	// URL - The path to the blob.
	URL *string `json:"url,omitempty"`
	// Sequencer - An opaque string value representing the logical sequence of events for any particular blob name. Users can use standard string comparison to understand the relative sequence of two events on the same blob name.
	Sequencer *string `json:"sequencer,omitempty"`
	// StorageDiagnostics - For service use only. Diagnostic data occasionally included by the Azure Storage service. This property should be ignored by event consumers.
	StorageDiagnostics interface{} `json:"storageDiagnostics,omitempty"`
}

// SubscriptionDeletedEventData schema of the Data property of an EventGridEvent for a
// Microsoft.EventGrid.SubscriptionDeletedEvent.
type SubscriptionDeletedEventData struct {
	// EventSubscriptionID - The Azure resource ID of the deleted event subscription.
	EventSubscriptionID *string `json:"eventSubscriptionId,omitempty"`
}

// SubscriptionValidationEventData schema of the Data property of an EventGridEvent for a
// Microsoft.EventGrid.SubscriptionValidationEvent.
type SubscriptionValidationEventData struct {
	// ValidationCode - The validation code sent by Azure Event Grid to validate an event subscription. To complete the validation handshake, the subscriber must either respond with this validation code as part of the validation response, or perform a GET request on the validationUrl (available starting version 2018-05-01-preview).
	ValidationCode *string `json:"validationCode,omitempty"`
	// ValidationURL - The validation URL sent by Azure Event Grid (available starting version 2018-05-01-preview). To complete the validation handshake, the subscriber must either respond with the validationCode as part of the validation response, or perform a GET request on the validationUrl (available starting version 2018-05-01-preview).
	ValidationURL *string `json:"validationUrl,omitempty"`
}

// SubscriptionValidationResponse to complete an event subscription validation handshake, a subscriber can use
// either the validationCode or the validationUrl received in a SubscriptionValidationEvent. When the
// validationCode is used, the SubscriptionValidationResponse can be used to build the response.
type SubscriptionValidationResponse struct {
	// ValidationResponse - The validation response sent by the subscriber to Azure Event Grid to complete the validation of an event subscription.
	ValidationResponse *string `json:"validationResponse,omitempty"`
}