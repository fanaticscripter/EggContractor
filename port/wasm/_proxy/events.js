// EGG_CONTRACTOR_BASEURL is an environment variable set by the CF Worker. It
// should be the base URL of a deployed EggContractor instance, e.g.
// "https://egg.your.domain".

const API_URL = `${EGG_CONTRACTOR_BASEURL}/api/events/`;

async function handleRequest(request) {
  const url = new URL(request.url);
  if (url.pathname !== "/") {
    return new Response(null, {
      status: 404,
      statusText: "Not Found",
    });
  }
  return await fetch(API_URL);
}

addEventListener("fetch", event => {
  return event.respondWith(handleRequest(event.request));
});
