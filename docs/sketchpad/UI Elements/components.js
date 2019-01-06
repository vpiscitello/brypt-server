let rootCluster = {
    coordinators: [{
        "id": 1,
        "cluster": 1,
        "coordinator": 0,
        "neighbor_count": 6,
        "designation": "root",
        "comm_techs": ["WiFi", "BLE", "LoRa"],
        "add_timestamp": 1545027824,
        "update_timestamp": 1545042224,
        "ireg_rate": 0.05
    }],
    neighbors: [{
        "id": 2,
        "cluster": 1,
        "coordinator": 1,
        "neighbor_count": 6,
        "designation": "node",
        "comm_techs": ["WiFi"],
        "add_timestamp": 1545027990,
        "update_timestamp": 1545042222,
        "ireg_rate": 0.1
    }, {
        "id": 3,
        "cluster": 1,
        "coordinator": 1,
        "neighbor_count": 6,
        "designation": "node",
        "comm_techs": ["WiFi", "BLE"],
        "add_timestamp": 1545027991,
        "update_timestamp": 1545042222,
        "ireg_rate": 0.12
    }, {
        "id": 4,
        "cluster": 1,
        "coordinator": 1,
        "neighbor_count": 9,
        "designation": "coordinator",
        "comm_techs": ["WiFi", "BLE", "LoRa"],
        "add_timestamp": 1545027992,
        "update_timestamp": 1545042223,
        "ireg_rate": 0.05
    }, {
        "id": 5,
        "cluster": 1,
        "coordinator": 1,
        "neighbor_count": 6,
        "designation": "node",
        "comm_techs": ["BLE"],
        "add_timestamp": 1545027993,
        "update_timestamp": 1545042222,
        "ireg_rate": 0.24
    }, {
        "id": 6,
        "cluster": 1,
        "coordinator": 1,
        "neighbor_count": 6,
        "designation": "node",
        "comm_techs": ["LoRa"],
        "add_timestamp": 1545027994,
        "update_timestamp": 1545042222,
        "ireg_rate": 0.55
    }]
};
let subCluster = {
    coordinators: [{
        "id": 1,
        "cluster": 1,
        "coordinator": 0,
        "neighbor_count": 6,
        "designation": "root",
        "comm_techs": ["WiFi", "BLE", "LoRa"],
        "add_timestamp": 1545027824,
        "update_timestamp": 1545042222,
        "ireg_rate": 0.05
    }, {
        "id": 4,
        "cluster": 1,
        "coordinator": 1,
        "neighbor_count": 9,
        "designation": "coordinator",
        "comm_techs": ["WiFi", "BLE", "LoRa"],
        "add_timestamp": 1545027992,
        "update_timestamp": 1545042223,
        "ireg_rate": 0.05
    }],
    neighbors: [{
        "id": 7,
        "cluster": 2,
        "coordinator": 4,
        "neighbor_count": 4,
        "designation": "node",
        "comm_techs": ["BLE"],
        "add_timestamp": 1545293042,
        "update_timestamp": 1545042223,
        "ireg_rate": 0.06
    }, {
        "id": 8,
        "cluster": 2,
        "coordinator": 4,
        "neighbor_count": 4,
        "designation": "node",
        "comm_techs": ["BLE"],
        "add_timestamp": 1545293043,
        "update_timestamp": 1545042223,
        "ireg_rate": 0.08
    }, {
        "id": 9,
        "cluster": 2,
        "coordinator": 4,
        "neighbor_count": 4,
        "designation": "node",
        "comm_techs": ["BLE"],
        "add_timestamp": 1545293044,
        "update_timestamp": 1545042223,
        "ireg_rate": 0.05
    }]
};


function indexOfObject(obj, arr) {
    for (let idx = 0; idx < arr.length; idx++) {
        if (obj.id == arr[idx].id) {
            return idx;
        }
    }
    return -1;
}

Vue.use(Vuex);

let store = new Vuex.Store({
    state: {
        connected: false,
        error: '',
        message: '',
        data: {
            aggReading: null,
            avgReading: null,
            nodeCount: null,
            clusterCount: null,
            attackCount: null,
            uptime: null
        }
    },
    mutations: {
        SOCKET_CONNECT(state) {
                state.connected = true;
            },
            SOCKET_DISCONNECT(state) {
                state.connected = false;
            },
            SOCKET_MESSAGE(state, message) {
                state.message = message;
            },
            SOCKET_ERROR(state, message) {
                state.error = message.error;
            }
    }
});


// // https://stackoverflow.com/questions/36170425/detect-click-outside-element
Vue.directive('click-outside', {
    bind: function(el, binding, vnode) {
        el.clickOutsideEvent = function(event) {
            if (!(el == event.target || el.contains(event.target))) {
                vnode.context[binding.expression](event);
            }
        };
        document.body.addEventListener('click', el.clickOutsideEvent);
    },
    unbind: function(el) {
        document.body.removeEventListener('click', el.clickOutsideEvent);
    }
});

// TODO: Get Network connection from Central Server
// TODO: Set up Control and Data websockets to network

let Spinner = {
    template: '#spinner-template',
    props: {
        // The color of the spinner elements
        color: {
            type: String,
            required: true
        }
    }
};


let FlagMessage = {
    template: '#flag-template',
    props: {
        message: {
            type: String,
            required: true
        },
        urgency: {
            type: Number,
            default: 0
        }
    },
    data: function() {
        return {
            show: true
        };
    },
    computed: {
        alertClass: function() {
            let alertClass = "";
            switch (this.urgency) {
                case 2:
                    alertClass = "danger";
                    // color = "#E74C3C";
                    break;
                case 1:
                    alertClass = "warning";
                    // color = "#FFB75E";
                    break;
                default:
                    alertClass = "basic";
                    // color = "#19CC95";
                    break;
            }
            return alertClass;
        }
    },
    methods: {
        showFlash: function() {
            this.show = true;
        },
        closeFlash: function() {
            this.show = false;
        }
    },
    created: function() {

    }
};


// Based on https://github.com/johndatserakis/vue-simple-context-menu
let ItemContextMenu = {
    template: '#ctx-menu-template',
    props: {
        // The element id of the context menu
        id: {
            type: String,
            required: true
        },
        // The supplied array of possible actions
        options: {
            type: Array,
            required: true
        }
    },
    data: function() {
        return {
            item: null, // The item that the context menu is operating on
            menuWidth: null, // The calculated width of the menu
            menuHeight: null, // The calculated height of the menu
            itemIsBranching: false // A boolean value set when the item is a coordinator
        };
    },
    computed: {
        // The options available to a leaf node
        leafOptions: function() {
            // Filter out options that need a branching node
            return this.options.filter((opt) => {
                return !opt.requiresCordinator;
            });
        },
        // The options available to a branch node
        branchOptions: function() {
            return this.options;
        }
    },
    methods: {
        showMenu: function(event, item) {
            // console.log("Show Menu");

            this.item = item; // Store the item being operated on

            // Set and Update the context menu based on the designation of the item
            this.itemIsBranching = (this.item.designation == "node") ? false : true;

            // After the DOM has been updated display the context menu
            this.$nextTick(() => {
                // Select the context menu and return if it does not exist
                let menu = document.getElementById(this.id);
                if (!menu) {
                    return;
                }

                // Compute the width of the menu if it has not been set
                if (!this.menuWidth || !this.menuHeight) {
                    menu.style.visibility = "hidden";
                    menu.style.display = "block";
                    this.menuWidth = menu.offsetWidth;
                    this.menuHeight = menu.offsetHeight;
                    menu.removeAttribute("style");
                }

                let parentRect = this.$parent.$el.getBoundingClientRect();

                // Set the X postion of the menu based on the item and page position
                if ((this.menuWidth + event.pageX) >= parentRect.width) {
                    menu.style.left = (event.pageX - parentRect.x - this.menuWidth + 2) + "px";
                } else {
                    menu.style.left = (event.pageX - parentRect.x - 2) + "px";
                }

                // Set the Y postion of the menu based on the item and page position
                if ((this.menuHeight + event.pageY) >= window.innerHeight) {
                    menu.style.top = (event.pageY - parentRect.y - this.menuHeight + 2) + "px";
                } else {
                    menu.style.top = (event.pageY - parentRect.y - 2) + "px";
                }

                menu.classList.add('active'); // Display the menu
            });
        },
        hideMenu: function() {
            // console.log("Hide Menu");
            let element = document.getElementById(this.id); // Get the menu
            if (element) {
                element.classList.remove('active'); // Hide the menu
            }
        },
        optionClicked: function(option) {
            this.hideMenu();
            this.$emit(option.event, this.item); // Notify the parent of a selected option
        },
        onClickOutside: function() {
            this.hideMenu();
        }
    }
};

let NodeDetail = {
    template: '#detail-template',
    props: {
        id: {
            type: String,
            required: true
        }
    },
    data: function() {
        return {
            show: false,
            node: null // The node the detail modal is providing more information on
        };
    },
    computed: {
        // Defines whether or not the node exhibits abnormal behavior on the network
        irregular: function() {
            return this.node.ireg_rate > 0.25;
        }
    },
    filters: {
        toTimeString: function(timestamp) {
            let date = new Date(timestamp * 1000); // Convert the provided timestamp to a Date object
            // Return a Locale Date string in the format Month Day, Year, Hour:Minute
            return date.toLocaleString(window.navigator.language, {
                day: "numeric",
                month: "short",
                year: "numeric",
                hour: "numeric",
                minute: "2-digit",
                hour12: true
            });
        },
        arrayToString: function(arr) {
            return arr.join(", "); // Join an array seperated by ", " i.e. "1, 2, 3"
        },
        toPercentage: function(rate) {
            return parseFloat(rate * 100).toFixed(1) + "%";
        }
    },
    methods: {
        showModal: function(node) {
            // console.log("Show Detail");
            this.node = node; // Store the node being operated on
            this.show = true;
        },
        hideModal: function() {
            // console.log("Hide Detail");
            this.show = false;
        },
        removeDevice: function() {
            console.log("Remove Node ", this.node.id, "From the User's Network.");
            // Launch Confirmation Modal
            // Send Remove Request for the node to the Network
            // Fetch the updated cluster
            alert("Sorry, this functionality is not yet available!");
        },
        editDevice: function() {
            console.log("Edit the details of Node ", this.node.id);
            // Launch Edit Modal
            alert("Sorry, this functionality is not yet available!");
        },
        toggleShow: function() {
            this.show = !this.show;
        },
        onClickOutside: function() {
            this.hideModal();
        }
    }
};

let NodeItem = {
    template: '#node-template',
    props: {
        node: {
            type: Object,
            required: true,
            validator: function(node) {
                return (
                    typeof node.id === "number" &&
                    typeof node.cluster === "number" &&
                    typeof node.coordinator === "number" &&
                    typeof node.neighbor_count === "number" &&
                    typeof node.designation === "string" &&
                    typeof node.comm_techs === "object" &&
                    typeof node.add_timestamp === "number" &&
                    typeof node.update_timestamp === "number" &&
                    typeof node.ireg_rate === "number"
                );
            }
        }
    },
    data: function() {
        return {

        };
    },
    computed: {
        // Defines whether or not the node exhibits abnormal behavior on the network
        irregular: function() {
            return this.node.ireg_rate > 0.25;
        },
        // The left zero pading for the node ID
        zeroPad: function() {
            zeros = "0000";
            // Make a subtring the ensure all IDs are aligned on four characters
            return zeros.substring(0, zeros.length - ("" + this.node.id).length);
        }
    },
    filters: {
        toTimeString: function(timestamp) {
            let date = new Date(timestamp * 1000); // Convert the provided timestamp to a Date object
            // Return a Locale Date string in the format Month Day, Year, Hour:Minute
            return date.toLocaleString(window.navigator.language, {
                day: "numeric",
                month: "short",
                year: "numeric",
                hour: "numeric",
                minute: "2-digit",
                hour12: true
            });
        },
        arrayToString: function(arr) {
            return arr.join(", "); // Join an array seperated by ", " i.e. "1, 2, 3"
        },
        capitalize: function(str) {
            return str.charAt(0).toUpperCase() + str.slice(1); // Capitalize the first letter of a string
        }
    },
    methods: {
        handleClick: function(event, node) {
            // Send an event to the root listener about a click on the item
            this.$root.$emit("item-action-menu-triggered", {
                event: event,
                node: node
            });
        }
    }
};

let NodeContainer = {
    template: '#ng-template',
    components: {
        'spinner': Spinner,
        'node-item': NodeItem
    },
    props: {
        id: {
            type: String,
            required: true
        },
        // The title of the node group
        title: {
            type: String,
            required: true
        },
        // The array of nodes in the group
        nodes: {
            type: Array,
            required: true
        },
        // Defines whether or not the group needs an async loading spinner
        needSpinner: {
            type: Boolean,
            required: true
        }
    },
    data: function() {
        return {
            ui: {
                loading: true // A boolean value tracking the updating state of the group
            }
        };
    },
    computed: {
        // The total height of items in the group
        groupHeight: function() {
            return (60 * this.nodes.length);
        },
        containerHeight: function() {
            return (this.groupHeight + 30);
        }
    },
    methods: {
        toggleLoader: function() {
            this.ui.loading = !this.ui.loading; // Toggle the loading tracker for the group
        }
    }
};

let ClusterContext = {
    template: '#cluster-template',
    components: {
        'spinner': Spinner,
        'node-container': NodeContainer,
        'item-context-menu': ItemContextMenu,
        'node-detail': NodeDetail
    },
    props: {

    },
    data: function() {
        return {
            coordinators: null, // The array of coordinators for the cluster
            neighbors: null // The array of neighbors to the cluster coordinator
        };
    },
    computed: {
        // Possible action menu options for the clusters
        actionMenuOptions: function() {
            let options = [];
            // option: {
            //     name: String,    // The name of the option for the user
            //     call: String,    // The function to be called upon selecting the option
            //     event: String,   // The event to trigger upon selecting the option
            //     requiresCordinator: Boolean  // A boolean value for whether or not the option needs a branching node
            // }
            // Push an option to open the node detail modal
            options.push({
                name: "Open Detail",
                call: "openDetail",
                event: "open-detail",
                requiresCordinator: false
            });
            // Push an option to expand (fetch) the neighbors of a coordinator
            options.push({
                name: "Expand Cluster",
                call: "expandCluster",
                event: "expand-cluster",
                requiresCordinator: true
            });
            return options;
        }
    },
    methods: {
        pushClusterCoordinator: function(node) {

            // Ensure a node to push has been provided
            if (node === null) {
                return;
            }

            let coordinators = []; // Intialize an empty array to store the coordinators
            let coordinatorStore = sessionStorage.getItem("coordinators"); // Read SessionStorage for stored coordinators

            // If there is comething in the Store parse the contents
            if (typeof coordinatorStore === "string") {
                coordinators = JSON.parse(coordinatorStore);
            }

            // Check if the node being pushed is in the array
            if (indexOfObject(node, coordinators) == -1) {
                coordinators.push(node); // Push the node into the coordinators array
                sessionStorage.setItem("coordinators", JSON.stringify(coordinators)); // Store the updated coordinators in SessionStorage
            }

            this.coordinators = coordinators; // Update the components data with the coordinator data

        },
        popClusterCoordinator: function() {
            this.coordinators.pop(); // Pop the last coordinator in the array
            sessionStorage.setItem("coordinators", JSON.stringify(this.coordinators)); // Store the updated coordinators in SessionStorage
        },
        fetchClusterTable: function(node) {

            // Ensure a node has been provided
            if (node === null) {
                return;
            }

            // Ensure the node is a branch
            if (node.designation == "node") {
                return;
            }

            let neighbors = Object;

            // Temp for the hardcoded cluster
            switch (node.id) {
                case 1:
                    {
                        console.log("Root Cluster");
                        neighbors = rootCluster.neighbors;
                        break;
                    }
                case 4:
                    {
                        console.log("SubCluster Query");
                        neighbors = subCluster.neighbors;
                        break;
                    }
                default:
                    {
                        console.log("No Cluster Exists for Coordinator");
                        break;
                    }
            }

            sessionStorage.setItem("neighbors", JSON.stringify(neighbors)); // Store the the neighbors in the SessionStorage in case of page reload
            this.neighbors = neighbors;

        },
        expandCluster: function(node) {
            console.log("Expand Cluster using: ", node.id);

            // Toggle loading state for coordinators and neighbors
            this.$refs.coordinators.toggleLoader();
            this.$refs.neighbors.toggleLoader();

            cordIndex = indexOfObject(node, this.coordinators); // Attempt to find the node in the array of coordinators
            // If the node is not in the coordinators push it
            // Otherwise pop all nodes up to its location
            if (cordIndex == -1) {
                this.pushClusterCoordinator(node);
            } else {
                // Pop each node until the request coordinator
                for (let idx = this.coordinators.length - 1; idx > cordIndex; idx--) {
                    this.popClusterCoordinator();
                }
            }
            this.$nextTick(() => {
                this.$refs.coordinators.toggleLoader(); // Toggle the loading state of the coordinators
            });

            this.fetchClusterTable(node); // Fetch and update neighbors data
            // this.$nextTick(() => {
            //     this.$refs.neighbors.toggleLoader(); // Toggle the loading state of the neighbors to close the spinner
            // });
            // Temp timeout to mimic async data loading
            setTimeout(function(t) {
                t.$nextTick(() => {
                    t.$refs.neighbors.toggleLoader();
                });
            }, 1500, this);
        },
        openDetail: function(node) {
            console.log("Open Detail using: ", node.id);
            this.$refs.detailModal.showModal(node); // Display the detail modal for the selected node
        },
        itemActionMenuTriggered: function(event, node) {
            if (this.$refs.detailModal.show) {
                this.$refs.detailModal.hideModal();
            }
            this.$refs.itemActionMenu.showMenu(event, node); // Display the context menu for the selected node
        }
    },
    created: function() {
        // Listen for clicks to trigger an action menu on a selected node
        this.$root.$on('item-action-menu-triggered', (args) => {
            this.itemActionMenuTriggered(args.event, args.node);
        });
        // Listen for events to expand a cluster
        this.$on('expand-cluster', (node) => {
            this.expandCluster(node);
        });
        // Listen for events to open a detail modal for a node
        this.$on('open-detail', (node) => {
            this.openDetail(node);
        });
    },
    beforeMount: function() {
        // TODO: Pass Root Cluster query

        // TEMP hardcoded root
        let root = {
            "id": 1,
            "cluster": 1,
            "coordinator": 0,
            "neighbor_count": 6,
            "designation": "root",
            "comm_techs": ["WiFi", "BLE", "LoRa"],
            "add_timestamp": 1545027824,
            "update_timestamp": 1545042224,
            "ireg_rate": 0.05
        };

        let coordinatorStore = sessionStorage.getItem("coordinators"); // Read SessionStorage for stored coordinators
        // If the coordinatorStore is empty push the root node
        // Otherwise load the stored information. This handles the case of page reload
        if (coordinatorStore === null) {
            this.pushClusterCoordinator(root);
        } else {
            this.coordinators = JSON.parse(coordinatorStore);
        }

        let neighborStore = sessionStorage.getItem("neighbors"); // Read SessionStorage for stored neighbors
        // If the neighborStore is empty query the network for the root cluster table
        // Otherwise load the stored information. This handles the case of page reload
        if (neighborStore === null) {
            this.fetchClusterTable(root);
        } else {
            this.neighbors = JSON.parse(neighborStore);
        }

    },
    mounted: function() {
        this.$nextTick(() => {
            this.$refs.coordinators.toggleLoader();
            // this.$refs.neighbors.toggleLoader();
        });
        // Temp timeout to mimic async data loading
        setTimeout(function(t) {
            t.$nextTick(() => {
                t.$refs.neighbors.toggleLoader();
            });
        }, 1500, this);
    },
    updated: function() {

    }
};

function CustomTooltip(tooltipModel) {
    var tooltipEl = document.getElementById('chartjs-tooltip');

    if (!tooltipEl) {
        tooltipEl = document.createElement('div');
        tooltipEl.id = 'chartjs-tooltip';
        tooltipEl.innerHTML = "<table></table>";
        document.body.appendChild(tooltipEl);
    }

    if (tooltipModel.opacity === 0) {
        tooltipEl.style.opacity = 0;
        return;
    }

    tooltipEl.classList.remove('above', 'below', 'no-transform');
    if (tooltipModel.yAlign) {
        tooltipEl.classList.add(tooltipModel.yAlign);
    } else {
        tooltipEl.classList.add('no-transform');
    }

    function getBody(bodyItem) {
        return bodyItem.lines;
    }

    if (tooltipModel.body) {
        var titleLines = tooltipModel.title || [];
        var bodyLines = tooltipModel.body.map(getBody);

        var innerHtml = '<tbody>';
        titleLines.forEach(function(title) {
            innerHtml += '<tr><th>' + title + '</th></tr>';
        });
        innerHtml += '</tbody>';

        var tableRoot = tooltipEl.querySelector('table');
        tableRoot.innerHTML = innerHtml;
    }

    var position = this._chart.canvas.getBoundingClientRect();


    tooltipEl.style.opacity = 1;
    tooltipEl.style.position = 'absolute';
    tooltipEl.style.left = position.left + window.pageXOffset + tooltipModel.caretX - (tooltipEl.offsetWidth /
            2) +
        'px';
    tooltipEl.style.top = position.bottom + window.pageYOffset - 50 + 'px';
    tooltipEl.style.color = '#FBFBFB';
    tooltipEl.style.fontFamily = tooltipModel._bodyFontFamily;
    tooltipEl.style.fontSize = tooltipModel.bodyFontSize + 'px';
    tooltipEl.style.fontStyle = tooltipModel._bodyFontStyle;
    tooltipEl.style.padding = 4 + 'px ' + 4 + 'px';
    tooltipEl.style.pointerEvents = 'none';
    tooltipEl.style.backgroundColor = '#1E1E1E';
    tooltipEl.style.borderRadius = 5 + 'px';
}

Chart.defaults.LineWithLine = Chart.defaults.line;

Chart.controllers.LineWithLine = Chart.controllers.line.extend({
    draw: function(ease) {
        Chart.controllers.line.prototype.draw.call(this, ease);

        if (this.chart.tooltip._active && this.chart.tooltip._active.length) {
            let activePoint = this.chart.tooltip._active[0],
                ctx = this.chart.ctx,
                x = activePoint.tooltipPosition().x,
                topY = this.chart.scales['y-axis-0'].top,
                bottomY = this.chart.scales['y-axis-0'].bottom;

            ctx.save();
            ctx.beginPath();
            ctx.moveTo(x, topY);
            ctx.lineTo(x, bottomY);
            ctx.lineWidth = 2;
            ctx.strokeStyle = 'rgba(251, 251, 251, 0.2)';
            ctx.stroke();
            ctx.restore();
        }
    }
});

const LineWithLine = VueChartJs.generateChart('custom-line', 'LineWithLine');

let ChartContainer = {
    // template: '#chart-template',
    extends: LineWithLine,
    mixins: [VueChartJs.mixins.reactiveProp],
    components: {
        'spinner': Spinner,
    },
    data: function() {
        return {
            gradients: null,
            options: {
                reactive: true,
                responsive: true,
                maintainAspectRatio: false,
                defaultFontFamily: Chart.defaults.global.defaultFontFamily = 'Source Sans Pro',
                scales: {
                    yAxes: [{
                        ticks: {
                            display: true,
                            fontColor: '#60777F',
                            fontSize: 12,
                        },
                        gridLines: {
                            display: true,
                            color: 'rgba(96, 119, 128, 0.1)'
                        },
                        scaleLabel: {
                            display: false,
                            labelString: 'Reading',
                            fontColor: '#60777F',
                            fontSize: 14
                        }
                    }],
                    xAxes: [{
                        type: 'realtime',
                        time: {
                            displayFormats: {
                                second: 'h:mm:ss A',
                            }
                        },
                        realtime: {
                            duration: 150000, // Show 10 30 second datapoints
                            // refresh: 30000, // refresh every 30 seconds
                            delay: 30000, // Delay of the show
                            onRefresh: function(chart) {

                            }

                        },
                        ticks: {
                            display: true,
                            fontColor: '#60777F',
                            fontSize: 12
                        },
                        gridLines: {
                            display: false
                        },
                        scaleLabel: {
                            display: true,
                            labelString: 'Timestamp',
                            fontColor: 'rgba(251, 251, 251, 0.6)',
                            fontSize: 13
                        }
                    }]
                },
                legend: {
                    display: true,
                    labels: {
                        boxWidth: 0,
                        fontColor: 'rgba(251, 251, 251, 0.6)',
                        // fontStyle: 'bold',
                        fontSize: 14
                    }
                },
                elements: {
                    line: {
                        tension: 0.2
                    }
                },
                hover: {
                    intersect: false
                },
                tooltips: {
                    mode: 'index',
                    intersect: false,
                    position: 'nearest',
                    titleFontSize: 0,
                    titleSpacing: 0,
                    titleMarginBottom: 0,
                    displayColors: false,
                    callbacks: {
                        title: function(tooltipItems, data) {
                            return tooltipItems[0].xLabel.toLocaleString(window.navigator.language, {
                                hour: "numeric",
                                minute: "2-digit",
                                second: "2-digit",
                                hour12: true
                            });
                        },
                        label: function(tooltipItem, data) {
                            return 'Reading: ' + tooltipItem.yLabel + '°';
                        }
                    },
                    custom: CustomTooltip
                },
                plugins: {
                    streaming: {
                        frameRate: 30
                    }
                }
            }
        };
    },
    methods: {
        getBackgroundGradient: function(idx) {
            return this.gradients[idx];
        }
    },
    watch: {
        chartData: function() {
            this.$data._chart.update();
        }
    },
    mounted: function() {
        this.gradients = [];
        let gradient = this.$refs.canvas.getContext('2d').createLinearGradient(0, 0, 0, 450);
        gradient.addColorStop(0, 'rgba(26, 204, 148, 0.6)');
        gradient.addColorStop(0.55, 'rgba(126, 238, 203, 0)');
        this.gradients.push(gradient);

        this.renderChart(this.chartData, this.options);
    },
    updated: function() {

    }
};

let NetworkOverview = {
    template: '#overview-template',
    components: {
        'spinner': Spinner,
    },
    props: {
        statistics: {
            type: Object,
            required: true
        }
    },
    data: function() {
        return {

        };
    },
    filters: {
        getTimePassString: function(timestamp) {
            let str = "";
            if (timestamp > 60 && timestamp < 3600) {
                str = parseFloat(timestamp / 60).toFixed(2) + " minutes";
            } else if (timestamp > 3600 && timestamp < 86400) {
                str = parseFloat(timestamp / 3600).toFixed(2) + " hours";
            } else if (timestamp > 86400 && timestamp < 604800) {
                str = parseFloat(timestamp / 86400).toFixed(2) + " days";
            } else if (timestamp > 604800) {
                str = parseFloat(timestamp / 604800).toFixed(2) + " weeks";
            } else {
                str = timestamp + "seconds";
            }
            return str;
        }
    }
};


let DataContext = {
    template: '#data-template',
    components: {
        'spinner': Spinner,
        'chart-container': ChartContainer,
        'network-overview': NetworkOverview,
    },
    props: {

    },
    data: function() {
        return {
            networkData: null,
            statistics: null,
            chartData: {
                labels: [],
                datasets: [{
                    label: 'Aggregate Readings',
                    pointColor: 'transparent',
                    pointStrokeColor: 'transparent',
                    pointHighlightFill: 'rgba(26, 204, 148, 1)',
                    pointHighlightStroke: 'rgba(126, 238, 203, 0.3)',
                    bezierCurve: true,
                    cubicInterpolationMode: 'monotone',
                    borderColor: 'rgb(126, 238, 203)',
                    borderWidth: 2,
                    pointRadius: 0,
                    pointBackgroundColor: 'rgba(251, 251, 251, 0)',
                    pointBorderColor: 'rgba(251, 251, 251, 0)',
                    pointHoverRadius: 5,
                    pointBorderWidth: 9,
                    pointHoverBackgroundColor: 'rgba(251, 251, 251, 1)',
                    pointHoverBorderColor: 'rgba(251, 251, 251, 0.2)',
                    backgroundColor: 'rgba(126, 238, 203, 0.2)',
                    data: []
                }]
            }
        };
    },
    computed: {

    },
    methods: {
        fetchNetworkData: function() {
            let networkData = {
                aggReading: null,
                avgReading: null,
                nodeCount: null,
                clusterCount: null,
                attackCount: null,
                uptime: null
            };

            // TODO: Fetch network data from the websocket

            // Create a new reading
            let reading = {
                timestamp: new Date(),
                data: (Math.random() * (72 - 68) + 68).toFixed(2)
            };

            networkData.aggReading = reading;
            networkData.avgReading = 68.34;
            networkData.nodeCount = 9;
            networkData.clusterCount = 2;
            networkData.attackCount = 15;
            networkData.uptime = (14 * 60 * 60) + (25 * 60);

            this.networkData = networkData;

            // Update the chart with the new network data
            this.$nextTick(() => {
                this.pushAggReadingToChart();
            });
        },
        pushAggReadingToChart: function() {
            // Create a new chartData object to update the chart
            // Needs a new object otherwise vue-chartjs will not register the change correctly
            let chartData = {
                labels: this.chartData.labels,
                datasets: [{
                    label: 'Aggregate Readings',
                    pointColor: 'transparent',
                    pointStrokeColor: 'transparent',
                    pointHighlightFill: 'rgba(26, 204, 148, 1)',
                    pointHighlightStroke: 'rgba(126, 238, 203, 0.3)',
                    bezierCurve: true,
                    cubicInterpolationMode: 'monotone',
                    borderColor: 'rgb(126, 238, 203)',
                    borderWidth: 2,
                    pointRadius: 0,
                    pointBackgroundColor: 'rgba(251, 251, 251, 0)',
                    pointBorderColor: 'rgba(251, 251, 251, 0)',
                    pointHoverRadius: 5,
                    pointBorderWidth: 9,
                    pointHoverBackgroundColor: 'rgba(251, 251, 251, 1)',
                    pointHoverBorderColor: 'rgba(251, 251, 251, 0.2)',
                    backgroundColor: this.$refs.chart.getBackgroundGradient(0),
                    data: this.chartData.datasets[0].data
                }]
            };

            // Push the new label and data, chartjs-plugin-streaming will remove the items as they leave the chart
            chartData.labels.push(this.networkData.aggReading.timestamp);
            chartData.datasets[0].data.push(this.networkData.aggReading.data);

            this.chartData = chartData; // Update the bound chartData
        },
        updateStatistics: function() {

        }
    },
    created: function() {
        this.fetchNetworkData(); // Get initial data about the network
        this.$nextTick(() => {
            // Initiate an interval to fetch new network data every thirty seconds
            this.dataUpdateInterval = setInterval(() => this.fetchNetworkData(), 30 * 1000);
        });
    }
};


// Bootstrap the App
let app = new Vue({
    el: '#app',
    store,
    components: {
        'flag-message': FlagMessage,
        'cluster-context': ClusterContext,
        'data-context': DataContext
    },
    data: {

    }
});
