import WizardContinueOrStart from './WizardContinueOrStart.vue';

describe('filteredDiscoveryTemplates', () => {
  it('includes the cVRP v4.0.1 template only for v4.0.1', () => {
    const cvrpV4 = { model: { discoveryModel: { name: 'cVRP v4.0.0' } } };
    const cvrpV4_0_1 = { model: { discoveryModel: { name: 'cVRP v4.0.1' } } };
    const discoveryTemplates = [cvrpV4, cvrpV4_0_1];
    const filter = WizardContinueOrStart.computed.filteredDiscoveryTemplates;

    expect(filter.call({ selectedVersion: 'v4.0.1', discoveryTemplates })).toEqual([cvrpV4_0_1]);
    expect(filter.call({ selectedVersion: 'v4.0.0', discoveryTemplates })).toEqual([]);
  });
});
