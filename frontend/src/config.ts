import { scaleOrdinal } from "d3-scale";
import { schemeCategory10 } from "d3-scale-chromatic";
import { FilterGroup, generateQueryStringFromFilterGroups } from "./data/FilterUtils";

export const colors = scaleOrdinal(schemeCategory10).range();

export const knownNodesFilter: FilterGroup[] = [
  [{ name: 'name', value: 'geth' }],
  [{ name: 'name', value: 'nethermind' }],
  [{ name: 'name', value: 'turbogeth' }],
  [{ name: 'name', value: 'turbo-geth' }],
  [{ name: 'name', value: 'erigon' }],
  [{ name: 'name', value: 'besu' }],
  [{ name: 'name', value: 'openethereum' }],
  [{ name: 'name', value: 'ethereum-js' }],

  [{ name: 'name', value: 'atlas' }],
  [{ name: 'name', value: 'besu' }],
  [{ name: 'name', value: 'bor' }],
  [{ name: 'name', value: 'coregeth' }],
  [{ name: 'name', value: 'efireal' }],
  [{ name: 'name', value: 'egem' }],
  [{ name: 'name', value: 'erigon' }],
  [{ name: 'name', value: 'eth2' }],
  [{ name: 'name', value: 'getd' }],
  [{ name: 'name', value: 'geth-ethercore' }],
  [{ name: 'name', value: 'gexp' }],
  [{ name: 'name', value: 'go-galaxy' }],
  [{ name: 'name', value: 'go-opera' }],
  [{ name: 'name', value: 'go-photon' }],
  [{ name: 'name', value: 'grails' }],
  [{ name: 'name', value: 'gubiq' }],
  [{ name: 'name', value: 'gvns' }],
  [{ name: 'name', value: 'na' }],
  [{ name: 'name', value: 'pirl' }],
  [{ name: 'name', value: 'q-qk_node' }],
  [{ name: 'name', value: 'quai' }],
  [{ name: 'name', value: 'ronin' }],
  [{ name: 'name', value: 'swarm' }],
  [{ name: 'name', value: 'thor' }],
  [{ name: 'name', value: 'wormholes' }],


]

export const knownNodesFilterString = generateQueryStringFromFilterGroups(knownNodesFilter)

export const LayoutEightPadding = [4, 4, 4, 8]
export const LayoutTwoColumn = ["repeat(1, 1fr)", "repeat(1, 1fr)", "repeat(1, 1fr)", "repeat(2, 1fr)"]
export const LayoutTwoColSpan = [1, 1, 1, 2]
