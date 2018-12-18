var coordinators = [{
    "id": "00000001",
    "cluster": 0,
    "coordinator": null,
    "neighbors": 5,
    "class": "root",
    "comm_techs": ["WiFi", "BLE", "LoRa"],
    "add_timestamp": 1545027824,
    "update_timestamp": 1545042222,
    "ireg_rate": 0.05
}];

var neighbors = [{
    "id": "00000001",
    "cluster": 0,
    "coordinator": "00000001",
    "neighbors": 5,
    "class": "node",
    "comm_techs": ["WiFi"],
    "add_timestamp": 1545027990,
    "update_timestamp": 1545042222,
    "ireg_rate": 0.1
}, {
    "id": "00000002",
    "cluster": 0,
    "coordinator": "00000001",
    "neighbors": 5,
    "class": "node",
    "comm_techs": ["WiFi", "BLE"],
    "add_timestamp": 1545027991,
    "update_timestamp": 1545042222,
    "ireg_rate": 0.12
}, {
    "id": "00000003",
    "cluster": 0,
    "coordinator": "00000001",
    "neighbors": 5,
    "class": "coordinator",
    "comm_techs": ["WiFi", "BLE", "LoRa"],
    "add_timestamp": 1545027992,
    "update_timestamp": 1545042222,
    "ireg_rate": 0.05
}, {
    "id": "00000004",
    "cluster": 0,
    "coordinator": "00000001",
    "neighbors": 5,
    "class": "node",
    "comm_techs": ["BLE"],
    "add_timestamp": 1545027993,
    "update_timestamp": 1545042222,
    "ireg_rate": 0.24
}, {
    "id": "00000005",
    "cluster": 0,
    "coordinator": "00000001",
    "neighbors": 5,
    "class": "node",
    "comm_techs": ["LoRa"],
    "add_timestamp": 1545027994,
    "update_timestamp": 1545042222,
    "ireg_rate": 0.55
}];

var NodeItem = {
    template: '#node-template',
    props: {
        node: {
            id: String,
            cluster: Number,
            coordinator: String,
            neighbors: Number,
            class: String,
            commTechs: Array,
            addTimestamp: Number,
            updateTimestamp: Number,
            ireg_rate: Number
        }
    },
    data: function() {
        return {
            id: this.node.id
        };
    }
};

var NodeGroup = {
    template: '#ng-template',
    components: {
        'node-item': NodeItem
    },
    props: {
        nodes: Array
    },
    data: function() {
        return {

        };
    },
    // computed: {
    //
    // },
    filters: {
        capitalize: function(str) {
            return str.charAt(0).toUpperCase() + str.slice(1);
        }
    },
    methods: {
        expandCluster: function() {

        },
        openDetail: function() {

        }
    }
};

var ClusterContext = {
    template: '#cluster-template',
    components: {
        'node-group': NodeGroup,
    },
    props: {
        coordinators: Array,
        neighbors: Array
    },
    // data: function() {
    //
    // }
};

// bootstrap the demo
var clusters = new Vue({
    el: '#clusters',
    components: {
        'cluster-context': ClusterContext,
        'node-group': NodeGroup,
        'node-item': NodeItem
    },
    data: {
        coordinators: coordinators,
        neighbors: neighbors
    }
});
