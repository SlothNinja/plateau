export function useUser(header, i) {
    return {
      id: header.userIds[i],
      name: header.userNames[i],
      emailHash: header.userEmailHashes[i],
      gravType: header.userGravTypes[i],
    }
}
