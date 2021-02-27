import { ei } from "./proto";

const messages = {
  commonlyInspected: [
    "EggIncFirstContactRequest",
    "EggIncFirstContactResponse",
    "GetPeriodicalsRequest",
    "PeriodicalsResponse",
    "MissionRequest",
    "MissionResponse",
    "CompleteMissionResponse",
    "ContractCoopStatusRequest",
    "ContractCoopStatusResponse",
    "Backup",
    "SaveBackupResponse",
  ],
  otherArtifactRequestResponse: [
    "ArtifactsConfigurationRequest",
    "ArtifactsConfigurationResponse",
    "CraftArtifactRequest",
    "CraftArtifactResponse",
    "ConsumeArtifactRequest",
    "ConsumeArtifactResponse",
    "SetArtifactRequest",
    "SetArtifactResponse",
  ],
  otherCoopRequestResponse: [
    "CreateCoopRequest",
    "CreateCoopResponse",
    "JoinCoopRequest",
    "JoinCoopResponse",
    "AutoJoinCoopRequest",
    "UpdateCoopPermissionsRequest",
    "UpdateCoopPermissionsResponse",
    "LeaveCoopRequest",
    "GiftPlayerCoopRequest",
    "KickPlayerCoopRequest",
    "ContractCoopStatusUpdateRequest",
    "ContractCoopStatusUpdateResponse",
  ],
  other: [],
};

const seen = [].concat(
  messages.commonlyInspected,
  messages.otherArtifactRequestResponse,
  messages.otherArtifactRequestResponse
);
for (const name in ei) {
  // Make sure we only pick up capitalized names just in case some lower case
  // helpers are introduced in the future.
  if ("ABCDEFGHIJKLMNOPQRSTUVWXYZ".includes(name[0])) {
    if (!seen.includes(name)) {
      messages.other.push(name);
    }
  }
}
messages.other.sort();

const messageGroups = [
  {
    label: "Commonly inspected",
    messages: messages.commonlyInspected,
  },
  {
    label: "Other artifact-related requests & responses",
    messages: messages.otherArtifactRequestResponse,
  },
  {
    label: "Other coop-related requests & responses",
    messages: messages.otherCoopRequestResponse,
  },
  {
    label: "Other",
    messages: messages.other,
  },
];

export { messageGroups };
