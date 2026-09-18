import WizardContinueOrStart from './WizardContinueOrStart.vue';

describe('filteredDiscoveryTemplates', () => {
  it('does not include cVRP templates on the front page', () => {
    const cvrpV4 = { model: { discoveryModel: { name: 'cVRP v4.0.0' } } };
    const cvrpV4_0_1 = { model: { discoveryModel: { name: 'cVRP v4.0.1' } } };
    const discoveryTemplates = [cvrpV4, cvrpV4_0_1];
    const filter = WizardContinueOrStart.computed.filteredDiscoveryTemplates;

    expect(filter.call({ selectedVersion: 'v4.0.1', discoveryTemplates })).toEqual([]);
    expect(filter.call({ selectedVersion: 'v4.0.0', discoveryTemplates })).toEqual([]);
  });
});
