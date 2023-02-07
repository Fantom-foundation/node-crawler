interface NamedCount {
	name: string;
	count: number;
}

export function appendOtherGroup(list: NamedCount[] | undefined): [NamedCount[], number] {
  let unknownItemCount = 0

  if (!list) {
    return [[], unknownItemCount];
  }
  
  const otherItem: NamedCount = {
    name: 'Other',
    count: 0
  };

  const filteredList: NamedCount[] = []
  list.forEach((item) => {
    if (item.name === 'tmp'  || item.name === '') {
      unknownItemCount +=  item.count
    } else if (filteredList.length > 10) {
      otherItem.count += item.count
    } else {
      filteredList.push(item)
    }
  })

  if (otherItem.count > 0) {
    filteredList.push(otherItem)
  }

  return [filteredList, unknownItemCount];
}

export function getClients(list: NamedCount[] | undefined): [NamedCount[], number] {
  let unknownItemCount = 0

  if (!list) {
    return [[], unknownItemCount];
  }

  const otherItem: NamedCount = {
    name: 'Other',
    count: 0
  };

  const filteredList: NamedCount[] = []
  list.forEach((item) => {
    if (item.name.includes('_')) {
      return
    }
    if (item.name === 'tmp'  || item.name === '') {
      unknownItemCount +=  item.count
    } else if (filteredList.length > 10) {
      otherItem.count += item.count
    } else {
      filteredList.push(item)
    }
  })

  if (otherItem.count > 0) {
    filteredList.push(otherItem)
  }

  return [filteredList, unknownItemCount];
}

export function getErrors(list: NamedCount[] | undefined): NamedCount[] {
  let unknownItemCount = 0
  if (!list) {
    return [];
  }

  const otherItem: NamedCount = {
    name: 'Other',
    count: 0
  };

  const filteredList: NamedCount[] = []
  list.forEach((item) => {
    if (item.name.includes('_')) {
      if (filteredList.length > 10) {
        otherItem.count += item.count
      } else {
        filteredList.push(item)
      }
    }
  })

  return filteredList
}