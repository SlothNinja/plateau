export function useCreator(header) {
  return {
    id: header.creatorId,
    name: header.creatorName,
    emailHash: header.creatorEmailHash,
    gravType: header.creatorGravType
  }
}
