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
  [{ name: 'name', value: 'gubiq' }],
  [{ name: 'name', value: 'mix-geth' }],
  [{ name: 'name', value: 'q-client' }],
  [{ name: 'name', value: 'qit' }],
  [{ name: 'name', value: 'qk_node' }],

  [{ name: 'name', value: 'readhellofailure_too_many_peers' }],
  [{ name: 'name', value: 'readhellofailure_expected_input_list_for_main.disconnect' }],
  [{ name: 'name', value: 'readhellofailure_eof' }],
  [{ name: 'name', value: 'readhellofailure_connection_reset_by_peer' }],
  [{ name: 'name', value: 'writehellofailure_connection_reset_by_peer' }],

  [{ name: 'name', value: 'readstatuserror_too_many_peers' }],
  [{ name: 'name', value: 'readstatuserror_37' }],
  [{ name: 'name', value: 'readstatuserror_33' }],
  [{ name: 'name', value: 'readstatuserror_expected_input_list_for_main.disconnect' }],
  [{ name: 'name', value: 'readstatuserror_disconnect_requested' }],
  [{ name: 'name', value: 'readstatuserror_eof' }],
  [{ name: 'name', value: 'readstatuserror_subprotocol_error' }],
  [{ name: 'name', value: 'readstatuserror_useless_peer' }],
  [{ name: 'name', value: 'readstatuserror_error_decoding_networkid' }],

  [{ name: 'name', value: 'couldnotdial_connection_refused' }],
  [{ name: 'name', value: 'couldnotdial_connection_reset_by_peer' }],
  [{ name: 'name', value: 'couldnotdial_eof' }],
  [{ name: 'name', value: 'couldnotdial_i/o_timeout' }],
  [{ name: 'name', value: 'couldnotdial_no_route_to_host' }],






]

export const knownNodesFilterString = generateQueryStringFromFilterGroups(knownNodesFilter)

export const LayoutEightPadding = [4, 4, 4, 8]
export const LayoutTwoColumn = ["repeat(1, 1fr)", "repeat(1, 1fr)", "repeat(1, 1fr)", "repeat(2, 1fr)"]
export const LayoutTwoColSpan = [1, 1, 1, 2]
