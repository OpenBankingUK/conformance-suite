<template>
  <div class="d-flex flex-row flex-fill">
    <div class="d-flex align-items-start">
      <div class="panel w-50" style="height:300px">
          <div class="panel-heading"><h5>Continue Test</h5></div>
      </div>
      <div class="panel w-50" style="height:300px">

        <div class="panel-heading"><h5>Continue Test</h5></div>
        <div class="panel-body">
          <p>Upload a signed and compatible report to view the results or to rerun.</p>
          <b-button class="mb-4">Choose file</b-button>
        </div>
      </div>
    </div>

    <div class="d-flex align-items-end">

      <div class="panel w-100">

          <div class="panel-heading"><h5>Or choose from a selection of discovery templates to begin.</h5></div>


          <div class="panel-body">
              <b-card-group deck>
        <DiscoveryTemplateCard
          v-for="(template, index) in discoveryTemplates"
          :key="index"
          :discoveryModel="template.model.discoveryModel"
          :image="template.image"
        />
      </b-card-group>
      
        </div>


          
      </div>

      


      


      


    </div>
  </div>

  <!-- <div class="align-items-start">



  <div class="row">

    
    <div class="col panel">Column</div>
    <div class="col panel">Column</div>

    <div class="w-100"></div>
    <div class="col">Column</div>
    <div class="col">Column</div>
  </div>


  <div class="flex-row">
      <div class="panel flex-column">
        <h2>Continue or Select Disvovery Template</h2>

        <p>Upload a signed and compatible report to view the results or to rerun.</p>
        <b-button class="mb-4">Choose file</b-button>
      </div>

      <div class="panel flex-column">
        a
      </div>



    </div>  

    <div class="panel w-100">
      <p>Or start a new test. Choose from a selection of discovery templates to begin.</p>

      <b-card-group deck>
        <DiscoveryTemplateCard
          v-for="(template, index) in discoveryTemplates"
          :key="index"
          :discoveryModel="template.model.discoveryModel"
          :image="template.image"
        />
      </b-card-group>
    </div>
  </div>-->
</template>

<script>
import { mapGetters } from "vuex";
import DiscoveryTemplateCard from "./DiscoveryTemplateCard.vue";

export default {
  name: "ContinueOrStart",
  components: {
    DiscoveryTemplateCard
  },
  computed: {
    ...mapGetters("config", ["discoveryTemplates"])
  }
};
</script>

<style>

</style>
