let updateMap = undefined;          // Global callbacks for updating map and charts
let oldWorld = undefined;           // Keep a global variable of world topology to not repeat drawing
let oldPortals = undefined;
let oldCitiesMap = new Map();
let refreshHandler = undefined;
let refreshTimeout = 5000;
let refreshLast = new Date();

const buildPortalsForm = (portals) => {
    const div = document.getElementById("tc_left");

    portals.forEach(portal => {
        const form = document.createElement("form");
        form.id = portal.name;
        form.className = 'left-form';
        errorContent = portal.status.error ? ' - Failed to connect' : '';
        const ratioPerSec = portal.settings.request_ratio > 0 ? (portal.settings.request_ratio / 100) * 1.33 : 0.06;
        form.innerHTML = `
            <h3>${portal.name} <span id="${portal.name}_error" class="error">${errorContent}</span></h3>
            <div class="left-form-group">
                <label>Ratio</label>
                <input class="requests-total" type="range" min="0" max="100" step="1" value="${portal.settings.request_ratio}" onchange="updatePortalForm(this.id, this.value)" id="${portal.name}_ratio" />
                <small id="${portal.name}_ratio_label">${ratioPerSec.toFixed(2)} req/s</small>
            </div>
            <div class="left-form-group">
                <label>Device</label>                
                <input class="device-mobile" type="range" min="0" max="100" step="1" value="${portal.settings.devices.mobile}" onchange="updatePortalForm(this.id, this.value)" id="${portal.name}_device_mobile" />
                <small id="${portal.name}_device_mobile_label">mob: ${portal.settings.devices.mobile}%</small>
                <input class="device-web" type="range" min="0" max="100" step="1" value="${portal.settings.devices.web}" onchange="updatePortalForm(this.id, this.value)" id="${portal.name}_device_web" />
                <small id="${portal.name}_device_web_label">web: ${portal.settings.devices.web}%</small>                
            </div>
            <div class="left-form-group">
                <label>User</label>                
                <input class="user-new" type="range" min="0" max="100" step="1" value="${portal.settings.users.new}" onchange="updatePortalForm(this.id, this.value)" id="${portal.name}_user_new" />
                <small id="${portal.name}_user_new_label">new: ${portal.settings.users.new}%</small>
                <input class="user-registered" type="range" min="0" max="100" step="1" value="${portal.settings.users.registered}" onchange="updatePortalForm(this.id, this.value)" id="${portal.name}_user_registered" />
                <small id="${portal.name}_user_registered_label">reg: ${portal.settings.users.registered}%</small>                
            </div>
            <div class="left-form-group">
                <label>Travel</label>                
                <input class="travel-t1" type="range" min="0" max="100" step="1" value="${portal.settings.travel_type.t1}" onchange="updatePortalForm(this.id, this.value)" id="${portal.name}_travel_t1" />
                <small id="${portal.name}_travel_t1_label">t1: ${portal.settings.travel_type.t1}%</small>
                <input class="travel-t2" type="range" min="0" max="100" step="1" value="${portal.settings.travel_type.t2}" onchange="updatePortalForm(this.id, this.value)" id="${portal.name}_travel_t2" />
                <small id="${portal.name}_travel_t2_label">t2: ${portal.settings.travel_type.t2}%</small>
                <input class="travel-t3" type="range" min="0" max="100" step="1" value="${portal.settings.travel_type.t3}" onchange="updatePortalForm(this.id, this.value)" id="${portal.name}_travel_t3" />
                <small id="${portal.name}_travel_t3_label">t3: ${portal.settings.travel_type.t3}%</small>                            
            </div>
        `;
        div.append(form);
    });

};

const buildRefreshForm = () => {
    const div = document.getElementById("tc_left");
    const form = document.createElement("form");
    form.id = "refresh-control";
    form.className = 'left-form';
    form.innerHTML = `
        <h2>Dashboard Settings</h2>
        <label>Refresh</label>
        <input class="refresh-control" type="range" min="0" max="10000" step="100" value="${refreshTimeout}" onchange="updateRefresh(this.value)" id="refresh_control" />
        <small id="refresh_control_label">${refreshTimeout} ms</small>
        <small id="refresh_status_label">Last refresh: ${refreshLast.toLocaleString()}</small>        
    `;
    div.append(form);
};

const drawChartsTotal = (portals) => {
    const div = document.getElementById("tc_charts_total");
    const width = div.offsetWidth - 10;
    const height = div.offsetHeight - 10;

    const margin = ({top: 40, right: 30, bottom: 10, left: 100})

    const svg = d3.select("#tc_charts_total")
        .append("svg")
        .attr("id", "svg_charts_total")
        .attr("viewBox", [0, 0, width, height]);

    const y = d3.scaleBand()
        .domain(d3.range(portals.length))
        .rangeRound([margin.top, height - margin.bottom])
        .padding(0.1);

    const yAxis = g => g
        .attr("transform", `translate(${margin.left},0)`)
        .call(d3.axisLeft(y).tickFormat(i => portals[i].name).tickSizeOuter(0))
        .attr("font-size", 14);

    const x = d3.scaleLinear()
        .domain([0, d3.max(portals, portal => portal.status.requests.total)])
        .range([margin.left, width - margin.right]);

    const xAxis = g => g
        .attr("transform", `translate(0,${margin.top})`)
        .call(d3.axisTop(x).ticks(width / 80))
        .call(g => g.select(".domain").remove())
        .call(g => g.select(".tick:first-of-type text").clone()
            .attr("y", -20)
            .attr("x", -100)
            .attr("text-anchor", "start")
            .attr("font-size", 16)
            .text("Total Requests per Portal"));

    const format = x.tickFormat(20);

    svg.append("g")
        .selectAll("rect")
        .data(portals)
        .join("rect")
        .attr("x", x(0))
        .attr("y", (d, i) => y(i))
        .attr("width", d => x(d.status.requests.total) - x(0))
        .attr("fill", (_, i) => "var(--portal-color-" + (i + 1) + ")")
        .attr("height", y.bandwidth());

    svg.append("g")
        .attr("fill", "white")
        .attr("text-anchor", "end")
        .attr("font-family", "sans-serif")
        .attr("font-size", 14)
        .selectAll("text")
        .data(portals)
        .join("text")
        .attr("x", d => x(d.status.requests.total))
        .attr("y", (d, i) => y(i) + y.bandwidth() / 2)
        .attr("dy", "0.35em")
        .attr("dx", -4)
        .text(d => format(d.status.requests.total))
        .call(text => text.filter(d => x(d.status.requests.total) - x(0) < 20) // short bars
            .attr("dx", +4)
            .attr("fill", "black")
            .attr("text-anchor", "start"));

    svg.append("g")
        .call(xAxis);

    svg.append("g")
        .call(yAxis);
};

const drawChartsDetails = (portals, divId, svgId, keys, colors, dataMap, title) => {
    const div = document.getElementById(divId);

    const width = div.offsetWidth - 10,
        height = div.offsetHeight - 10;

    const svg = d3.select("#" + divId)
        .append("svg")
        .attr("id", svgId)
        .attr("viewBox", [0, 0, width, height]);

    const margin = {top: 10, right: 10, bottom: 30, left: 40};

    const color = d3.scaleOrdinal()
        .range(colors);

    const x0 = d3.scaleBand()
        .domain(portals.map(d => d.name))
        .rangeRound([margin.left, width - margin.right])
        .paddingInner(0.1);

    const x1 = d3.scaleBand()
        .domain(keys)
        .rangeRound([0, x0.bandwidth()])
        .padding(0.05)

    const xAxis = g => g
        .attr("transform", `translate(0,${height - margin.bottom})`)
        .call(d3.axisBottom(x0).tickSizeOuter(0))
        .attr("font-size", 14)
        .call(g => g.select(".domain").remove())

    // Note that we know that d.status.requests.total will represent max request per city
    const y = d3.scaleLinear()
        .domain([0, d3.max(portals, d => d.status.requests.total)]).nice()
        .rangeRound([height - margin.bottom, margin.top])

    const yAxis = g => g
        .attr("transform", `translate(${margin.left},0)`)
        .call(d3.axisLeft(y).ticks(null, "s"))
        .call(g => g.select(".domain").remove())
        .call(g => g.select(".tick:last-of-type text").clone()
            .attr("x", 3)
            .attr("text-anchor", "start")
            .attr("font-size", 16)
            .text(title)
        );

    const legend = svg => {
        const g = svg
            .attr("transform", `translate(${width},0)`)
            .attr("text-anchor", "end")
            .attr("font-family", "sans-serif")
            .attr("font-size", 16)
            .selectAll("g")
            .data(color.domain().slice())
            .join("g")
            .attr("transform", (d, i) => `translate(0,${i * 20})`);

        g.append("rect")
            .attr("x", -19)
            .attr("width", 19)
            .attr("height", 19)
            .attr("fill", color);

        g.append("text")
            .attr("x", -24)
            .attr("y", 9.5)
            .attr("dy", "0.35em")
            .text(d => d);
    };

    svg.append("g")
        .call(xAxis);

    svg.append("g")
        .call(yAxis);

    svg.append("g")
        .selectAll("g")
        .data(portals)
        .join("g")
        .attr("transform", d => `translate(${x0(d.name)},0)`)
        .selectAll("rect")
        .data(dataMap)
        .join("rect")
        .attr("x", d => x1(d.key))
        .attr("y", d => y(d.value))
        .attr("width", x1.bandwidth())
        .attr("height", d => y(0) - y(d.value))
        .attr("fill", d => color(d.key));

    svg.append("g")
        .call(legend);
};

const drawChartsDevices = (portals) => {
    const keys = ["mobile", "web"];
    const colors = ["var(--device-mobile-color)", "var(--device-web-color)"];
    const dataMap = d => keys.map(key => {
        let value = 0;
        switch (key) {
            case "web":
                value = d.status.requests.devices.web;
                break;
            case "mobile":
                value = d.status.requests.devices.mobile;
                break;
            default:
            // Nothing to do in default
        }
        return {
            key,
            value: value,
        }
    });
    drawChartsDetails(
        portals,
        "tc_charts_devices",
        "svg_charts_devices",
        keys,
        colors,
        dataMap,
        "Device Type"
    );
};

const drawChartsUsers = (portals) => {
    const keys = ["new", "registered"];
    const colors = ["var(--user-new-color)", "var(--user-registered-color)"];
    const dataMap = d => keys.map(key => {
        let value = 0;
        switch (key) {
            case "registered":
                value = d.status.requests.users.registered;
                break;
            case "new":
                value = d.status.requests.users.new;
                break;
            default:
            // Nothing to do in default
        }
        return {
            key,
            value: value,
        }
    });
    drawChartsDetails(
        portals,
        "tc_charts_users",
        "svg_charts_users",
        keys,
        colors,
        dataMap,
        "User Type"
    );
};

const drawChartsTypes = (portals) => {
    const keys = ["t1", "t2", "t3"];
    const colors = ["var(--travel-t1-color)", "var(--travel-t2-color)", "var(--travel-t3-color)"];
    const dataMap = d => keys.map(key => {
        let value = 0;
        switch (key) {
            case "t1":
                value = d.status.requests.travel_type.t1;
                break;
            case "t2":
                value = d.status.requests.travel_type.t2;
                break;
            case "t3":
                value = d.status.requests.travel_type.t3;
                break;
            default:
            // Nothing to do in default
        }
        return {
            key,
            value: value,
        }
    });
    drawChartsDetails(
        portals,
        "tc_charts_travel_types",
        "svg_charts_travel_types",
        keys,
        colors,
        dataMap,
        "Travel Type"
    );
};

const drawCharts = (portals) => {
    d3.select("#svg_charts_total").remove();
    d3.select("#svg_charts_devices").remove();
    d3.select("#svg_charts_users").remove();
    d3.select("#svg_charts_travel_types").remove();

    drawChartsTotal(portals);
    drawChartsDevices(portals);
    drawChartsUsers(portals);
    drawChartsTypes(portals);
};

const drawMap = (world, portals) => {
    const div = document.getElementById("tc_map");

    const width = div.offsetWidth - 10,    // Adjust the margin
        height = div.offsetHeight - 10;  // Adjust the margin

    const svg = d3.select("#tc_map")
        .append("svg")
        .attr("id", "svg_map")
        .attr("viewBox", [0, 0, width, height]);

    // Define map projection
    const projection = d3.geoMercator()
        .center([13, 52])                    // comment centrer la carte, longitude, latitude
        .translate([width / 2, height / 2])  // centrer l'image obtenue dans le svg
        .scale([width / 2.5]);                     // zoom, plus la valeur est petit plus le zoom est gros

    // Define path generator
    const path = d3.geoPath()
        .projection(projection);

    const countries = topojson.feature(world, world.objects.countries);
    const filtered = portals.map(p => p.country);
    svg.selectAll("path")
        .data(countries.features)
        .enter()
        .append("path")
        .attr("d", path)
        .attr("stroke", "var(--map-stroke-color)")
        .attr("fill", (d) => {
            if (filtered.includes(d.properties.name)) {
                return "var(--map-fill-selected-color)";
            }
            return "var(--map-fill-color)";
        });

    svg.append("g")
        .attr("id", "g_status_info");

    // Map Legend
    svg.append("text")
        .attr("font-size", 18)
        .attr("x", 5)
        .attr("y", 30)
        .text("Requests per City");

    svg.append("circle")
        .attr("cx", 200)
        .attr("cy", 25)
        .attr("r", 10)
        .attr("fill", "var(--total-color)");

    svg.append("text")
        .attr("font-size", 16)
        .attr("x", 220)
        .attr("y", 30)
        .text("total per city");

    svg.append("circle")
        .attr("cx", 350)
        .attr("cy", 25)
        .attr("r", 10)
        .attr("fill", "var(--new-city-color)");

    svg.append("text")
        .attr("font-size", 16)
        .attr("x", 370)
        .attr("y", 30)
        .text("new request");

    updateMap = (newPortals) => {
        const cities = getCities(newPortals);

        const radius = d3.scaleLinear()
            .domain([0, d3.max(cities.map(c => c.total))])
            .range([0, 10]);

        svg.select("#g_status_info").remove();
        svg.append("g")
            .attr("id", "g_status_info");

        const g = svg.select("#g_status_info")
            .selectAll("circle")
            .data(cities)
            .enter();

        g.append("circle")
            .attr("cx", p => projection(p.coordinates)[0])
            .attr("cy", p => projection(p.coordinates)[1])
            .attr("r", p => p.new ? radius(p.total) * 3 : radius(p.total))
            .attr("fill", p => p.new ? "var(--new-city-color)" : "var(--total-color)")
            .attr("fill-opacity", 1)
            .attr("stroke-opacity", 0)
            .on("mouseenter", (d, _) => {
                if (d3.select("#t" + d.city).empty()) {
                    const pieRadius = radius(d.total);
                    const pieX = projection(d.coordinates)[0];
                    const pieY = projection(d.coordinates)[1];

                    // Magic xOffset
                    const xOffset = ((d.city.length / 2) * 8) + 12;
                    const yOffset = pieRadius * 3;
                    svg.select("#g_status_info")
                        .append("text")
                        .attr("id", "t" + d.city)
                        .attr("x", projection(d.coordinates)[0] - xOffset)
                        .attr("y", projection(d.coordinates)[1] - yOffset)
                        .text(d.city + ' (' + d.total + ')');

                    const pieData = [];
                    Object.keys(d['totalPortals']).forEach((key, j) => {
                        pieData.push({
                            name: key,
                            value: d['totalPortals'][key],
                        });
                    });

                    const pie = d3.pie()
                        .sort(null)
                        .value(d => d.value);

                    const arc = d3.arc()
                        .innerRadius(pieRadius)
                        .outerRadius(pieRadius * 2)

                    const arcs = pie(pieData);

                    const arcLabel = d3.arc().innerRadius(pieRadius).outerRadius(pieRadius);

                    svg.select("#g_status_info")
                        .append("g")
                        .attr("id", "p" + d.city)
                        .attr("transform", "translate(" + pieX + "," + pieY + ")")
                        .attr("stroke", "white")
                        .selectAll("path")
                        .data(arcs)
                        .join("path")
                        .attr("fill", (_, i) => "var(--portal-color-" + (i + 1) + ")")
                        .attr("d", arc)
                        .append("title")
                        .text(d => `${d.data.name}: ${d.data.value.toLocaleString()}`);
                }
            })
            .on("mouseleave", (d, _) => {
                // Select text by id and then remove
                d3.select("#t" + d.city).remove();  // Remove text location
                d3.select("#p" + d.city).remove();  // Remove pie location
            })
            .transition()
            .attr("fill-opacity", 0)
            .attr("stroke-opacity", 1)
            .attr("r", p => radius(p.total))
            .delay(500)
            .duration(refreshTimeout);
    }
};

const getCities = (portals) => {
    const tempMap = new Map();
    portals.forEach(portal => {
        portal.status.cities.forEach(city => {
            let totalCity = tempMap.get(city.city);
            if (!totalCity) {
                totalCity = {
                    city: city.city,
                    coordinates: city.coordinates,
                    total: 0,
                    new: false,
                    totalPortals: {},
                }
            }
            totalCity['total'] = totalCity['total'] + city.requests.total;
            totalCity['totalPortals'][portal.name] = city.requests.total;
            tempMap.set(city.city, totalCity);
        });
    });

    const cities = [];
    tempMap.forEach(newCity => {
        const oldCity = oldCitiesMap.get(newCity.city);
        if (!oldCity || (oldCity && oldCity.total < newCity.total)) {
            newCity.new = true;
        }
        cities.push(newCity);
    });
    oldCitiesMap = tempMap;
    return cities;
};

const getPortalsStatus = () => {
    return fetch('status')
        .then(resp => {
            return resp.json();
        });
};

const updatePortals = () => {
    getPortalsStatus()
        .then(newPortals => {
            oldPortals = newPortals;
            updateMap(oldPortals);
            drawCharts(oldPortals);
            updatePortalsForm(oldPortals);
            refreshLast = new Date();
            updateRefreshStatus();
        });
};

const updatePortalForm = (formid, value) => {
    const index = formid.indexOf("_");
    const portal = formid.substring(0, index);
    const type = formid.substring(index + 1);
    console.info('Changing portal: ' + portal + ' type: ' + type + ' value: ' + value);
    for (let i = 0; i < oldPortals.length; i++) {
        if (oldPortals[i].name === portal) {
            switch (type) {
                case "ratio":
                    oldPortals[i].settings.request_ratio = parseInt(value);
                    break;
                case "device_mobile":
                    oldPortals[i].settings.devices.mobile = parseInt(value);
                    oldPortals[i].settings.devices.web = 100 - value;
                    document.getElementById(portal + "_device_web").value = oldPortals[i].settings.devices.web;
                    break;
                case "device_web":
                    oldPortals[i].settings.devices.web = parseInt(value);
                    oldPortals[i].settings.devices.mobile = 100 - value;
                    document.getElementById(portal + "_device_mobile").value = oldPortals[i].settings.devices.mobile;
                    break;
                case "user_new":
                    oldPortals[i].settings.users.new = parseInt(value);
                    oldPortals[i].settings.users.registered = 100 - value;
                    document.getElementById(portal + "_user_registered").value = oldPortals[i].settings.users.registered;
                    break;
                case "user_registered":
                    oldPortals[i].settings.users.registered = parseInt(value);
                    oldPortals[i].settings.users.new = 100 - value;
                    document.getElementById(portal + "_user_new").value = oldPortals[i].settings.users.new;
                    break;
                case "travel_t1":
                    oldPortals[i].settings.travel_type.t1 = parseInt(value);
                    oldPortals[i].settings.travel_type.t2 = (100 - value) / 2;
                    oldPortals[i].settings.travel_type.t3 = (100 - value) / 2;
                    document.getElementById(portal + "_travel_t2").value = oldPortals[i].settings.travel_type.t2;
                    document.getElementById(portal + "_travel_t3").value = oldPortals[i].settings.travel_type.t3;
                    break;
                case "travel_t2":
                    oldPortals[i].settings.travel_type.t2 = parseInt(value);
                    oldPortals[i].settings.travel_type.t1 = (100 - value) / 2;
                    oldPortals[i].settings.travel_type.t3 = (100 - value) / 2;
                    document.getElementById(portal + "_travel_t1").value = oldPortals[i].settings.travel_type.t1;
                    document.getElementById(portal + "_travel_t3").value = oldPortals[i].settings.travel_type.t3;
                    break;
                case "travel_t3":
                    oldPortals[i].settings.travel_type.t3 = parseInt(value);
                    oldPortals[i].settings.travel_type.t1 = (100 - value) / 2;
                    oldPortals[i].settings.travel_type.t2 = (100 - value) / 2;
                    document.getElementById(portal + "_travel_t1").value = oldPortals[i].settings.travel_type.t1;
                    document.getElementById(portal + "_travel_t2").value = oldPortals[i].settings.travel_type.t2;
                    break;
            }
            updateSettings(portal, oldPortals[i].settings);
        }
    }
};

const updatePortalsForm = (portals) => {
    portals.forEach(p => {
        document.getElementById(p.name + "_ratio").value = p.settings.request_ratio;
        const ratioPerSec = p.settings.request_ratio > 0 ? (p.settings.request_ratio / 100) * 1.33 : 0.06;
        document.getElementById(p.name + "_ratio_label").innerText = ratioPerSec.toFixed(2) + ' req/s';

        document.getElementById(p.name + "_device_mobile").value = p.settings.devices.mobile;
        document.getElementById(p.name + "_device_mobile_label").innerText = 'mob: ' + p.settings.devices.mobile + '%';
        document.getElementById(p.name + "_device_web").value = p.settings.devices.web;
        document.getElementById(p.name + "_device_web_label").innerText = 'web: ' + p.settings.devices.web + '%';

        document.getElementById(p.name + "_user_new").value = p.settings.users.new;
        document.getElementById(p.name + "_user_new_label").innerText = 'new: ' + p.settings.users.new + '%';
        document.getElementById(p.name + "_user_registered").value = p.settings.users.registered;
        document.getElementById(p.name + "_user_registered_label").innerText = 'reg: ' + p.settings.users.registered + '%';

        document.getElementById(p.name + "_travel_t1").value = p.settings.travel_type.t1;
        document.getElementById(p.name + "_travel_t1_label").innerText = 't1: ' + p.settings.travel_type.t1 + '%';
        document.getElementById(p.name + "_travel_t2").value = p.settings.travel_type.t2;
        document.getElementById(p.name + "_travel_t2_label").innerText = 't2: ' + p.settings.travel_type.t2 + '%';
        document.getElementById(p.name + "_travel_t3").value = p.settings.travel_type.t3;
        document.getElementById(p.name + "_travel_t3_label").innerText = 't3: ' + p.settings.travel_type.t3 + '%';

        span = document.getElementById(p.name + "_error");
        if (p.status.error) {
            span.innerText = ' - Failed to connect ';
        } else {
            span.innerText = '';
        }
    });
};

const updateRefresh = (value) => {
    console.info("Changing Global Refresh to " + value + " ms");
    refreshTimeout = value
    document.getElementById('refresh_control_label').innerText = `${refreshTimeout} ms`;
    clearInterval(refreshHandler);
    refreshHandler = window.setInterval(updatePortals, refreshTimeout);
};

const updateRefreshStatus = () => {
    const label = document.getElementById("refresh_status_label");
    label.innerText = "Last refresh: " + refreshLast.toLocaleString();
};

const updateSettings = (portal, settings) => {
    fetch('settings/' + portal, {
        method: 'PUT',
        body: JSON.stringify(settings),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(resp => {
        return resp.json();
    }).then(_ => {
        updatePortals();
    });
};

// Main entry points

Promise.all([
    d3.json("data/countries-110m.json"),    // From https://github.com/topojson/world-atlas
    getPortalsStatus(),
])
    .then(([newWorld, newPortals]) => {
        oldPortals = newPortals;
        oldWorld = newWorld;

        buildPortalsForm(oldPortals);
        buildRefreshForm();

        drawCharts(oldPortals);

        drawMap(oldWorld, oldPortals);
        updateMap(oldPortals);
    });

window.onresize = () => {
    if (oldWorld !== undefined && oldPortals !== undefined) {

        d3.select("#svg_map").remove();
        drawMap(oldWorld, oldPortals);
        updateMap(oldPortals);

        drawCharts(oldPortals);
    }
};

refreshHandler = window.setInterval(updatePortals, refreshTimeout);

