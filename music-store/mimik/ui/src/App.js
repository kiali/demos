import React, { useState } from 'react';

function App() {
    const [namespaces, setNamespaces] = useState(5);
    const [services, setServices] = useState(5);
    const [connections, setConnections] = useState(3);
    const [randomConnections, setRandomConnections] = useState(4);
    const [topology, setTopology] = useState({});
    const [error, setError] = useState("");

    const generateTopology = (e) => {
        e.preventDefault();
        console.log(JSON.stringify({
            "namespaces": namespaces,
            "services": services,
            "connections": connections,
            "randomConnections": randomConnections
        }));
        fetch("/generate", 
            { 
                method: 'post',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    "namespaces": namespaces,
                    "services": services,
                    "connections": connections,
                    "randomConnections": randomConnections
                })
            })
            .then(res => res.json())
            .then(
                (result) => {
                    console.log(result);
                    setTopology(result);
                },
                (error) => {
                    console.log(error);
                    setError("Error!");
                }
            )
    }

    return (
        <>
            <nav className="navbar navbar-expand-lg navbar-light bg-light">
                <div className="container">
                    <a className="navbar-brand" href="#">Topology Generator</a>
                </div>
            </nav>
            <div className="container" style={{ marginTop: "20px" }}>
                {error && <div className="row">
                    <div className="col">
                        <div className="alert alert-danger" role="alert">
                            {error}
                        </div>
                    </div>
                </div>}
                <div className="row">
                    <div className="col-4">
                        <div className="card ">
                            <div className="card-body">
                                <form>
                                    <label className="form-label">Namespaces</label>
                                    <input type="range" className="form-range" min="1" max="10" value={namespaces} id="namespacesRange" name="namespacesRange" onChange={(e) => { setNamespaces(parseInt(e.target.value)) }} />
                                    <p>{namespaces}</p>

                                    <label className="form-label">Services per namespace </label>
                                    <input type="range" className="form-range" min="1" max="20" value={services} id="servicesRange" name="servicesRange" onChange={(e) => { setServices(parseInt(e.target.value)) }} />
                                    <p>{services}</p>

                                    <label className="form-label">Connections per service</label>
                                    <input type="range" className="form-range" min="1" max="5" value={connections} id="connectionsRange" name="connectionsRange" onChange={(e) => { setConnections(parseInt(e.target.value)) }} />
                                    <p>{connections}</p>

                                    <label className="form-label">Random connections</label>
                                    <input type="range" className="form-range" min="1" max="10" value={randomConnections} id="randomConnectionsRange" name="randomConnectionsRange" onChange={(e) => { setRandomConnections(parseInt(e.target.value)) }} />
                                    <p>{randomConnections}</p>

                                    <button className="btn btn-primary float-end" onClick={generateTopology}>Generate</button>
                                </form>
                            </div>
                        </div>
                    </div>
                    <div className="col-8">
                        <textarea className="form-control" style={{ border: "1px solid rgba(0,0,0,.125)" }} rows="30" value={"cat <<EOF | kubectl create -f - \n" + JSON.stringify(topology, null, 2) + "\nEOF"} />
                    </div>

                </div>
                <div className="mb-3">

                </div>
            </div>
        </>
    );
}

export default App;
