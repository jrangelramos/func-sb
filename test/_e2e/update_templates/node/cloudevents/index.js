'use strict';
const { CloudEvent, HTTP } = require('cloudevents');

let lastInvokeEventData = ""

/**
 * Function used to test 'func invoke' command
 * The trick here is sending the event using emit with a given event source 'func:invoke'.
 * For this source the consumed event will be stored and returned as http response when it received
 * another event with source 'e2e:check'.
 *
 * 1) function will consume and store the data "hello invoke"
 * kn func invoke -c "text/plain" -d "hello invoke" -s "func:invoke"
 *
 * 2) the below should return "hello invoke" from previous command
 * curl $node_func_url -X POST \
 *  -H "Ce-Id: some-message-id" \
 *  -H "Ce-Specversion: 1.0" \
 *  -H "Ce-Type: e2e:check" \
 *  -H "Ce-Source: e2e:check" \
 *  -H "Content-Type: text/plain" \
 *  -d 'Invoke Check'
 *
 *
 * @param context
 * @param cloudevent
 * @returns {{message: string}|*}
 */
function handle(context, cloudevent) {
  if (!cloudevent) {
    return {
      message: 'No cloud event received'
    };
  }

  if (cloudevent.source == "func:invoke") {
    context.log.info(`CloudEvent received : ${cloudevent.toString()}`);
    lastInvokeEventData = cloudevent.data
  }

  if (cloudevent.source == "e2e:check") {
    return HTTP.binary(new CloudEvent({
      source: 'test:handle',
      type: 'test:invoke',
      data: lastInvokeEventData
    }));
  }

  return {
    message: 'Cloud event received'
  };
};

module.exports = handle;
