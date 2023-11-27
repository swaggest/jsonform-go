(function () {
    "use strict";

    /**
     * Form
     * @constructor
     */
    function JSONForm() {
        /**
         * @type {Element}
         */
        this.form = null;

        /**
         * @type {Element}
         */
        this.title = null;

        /**
         * @type {Element}
         */
        this.result = null;

        /**
         * @type {RawCallback}
         */
        this.fail = null;

        /**
         * @type {HTMLCallback}
         */
        this.error = null;

        /**
         * @type {RawCallback}
         */
        this.success = null;

        /**
         * @type {JSONCallback}
         */
        this.beforeSubmit = null

        this.schemaName = '';
        this.schema = undefined;

        this.valueUrl = '';
        this.value = undefined;

        this.submitUrl = '';
        this.submitMethod = 'POST';
        this.successStatus = 200;
    }

    /**
     * @typedef formParams
     * @type {Object}
     * @property {String} title - Title of the form.
     * @property {String} schemaName - Schema name.
     * @property {String} valueUrl - URL to fetch value.
     * @property {String} submitUrl - URL to submit form.
     * @property {String} submitMethod - HTTP method to use on form submit.
     * @property {Number} successStatus - Success HTTP status code to expect on submit.
     * @property {RawCallback} onSuccess - Callback for successful response.
     * @property {RawCallback} onFail - Callback for failed response.
     * @property {HTMLCallback} onError - Callback for error.
     * @property {JSONCallback} onBeforeSubmit - Callback for sumbittable form data.
     *
     * @property {Object} value - Value, can be absent if provided with valueUrl.
     * @property {Object} schema - Schema, can be absent if provided with schemaUrl.
     */


    JSONForm.prototype.default = function () {
        /**
         * @type {formParams}
         */
        var params = this.queryParams()

        this.make(params)
    }

    /**
     *
     * @param {formParams} params
     */
    JSONForm.prototype.make = function (params) {
        if (this.title === null) {
            this.title = $("#title");
        }

        if (this.result === null) {
            this.result = $("#res");
        }

        if (params.title) {
            $(this.title).text(params.title)
            if (document.title === "") {
                document.title = params.title
            }
        }
        console.log("QUERY PARAMS:", params)

        if (params.onSuccess) {
            this.success = params.onSuccess
        }

        if (params.onFail) {
            this.fail = params.onFail
        }

        if (params.onError) {
            this.error = params.onError
        }

        if (params.onBeforeSubmit) {
            this.beforeSubmit = params.onBeforeSubmit
        }

        var self = this

        if (this.error === null) {
            this.error = function (html) {
                this.result.html('ERROR: ' + html);
            }
        }

        if (this.fail === null) {
            this.fail = function (x) {
                self.error("Failed to submit form using URL:<br /><code>" + self.submitUrl + "</code><br />" +
                    "Expected status:<br /><code>" + self.successStatus + "</code><br />" +
                    "Status:<br /><code>" + x.status + "</code><br />" +
                    "Response:<br /><code>" + x.responseText + "</code>")
            }
        }

        if (this.success === null) {
            this.success = function (html) {
                self.result.html("Submitted.");
            }
        }

        if (params.schema != null) {
            this.schema = params.schema;
        } else {
            if (params.schemaName == null) {
                this.error("Missing schemaName parameter in URL");
                return;
            }
        }


        if (params.submitUrl == null) {
            this.error("Missing submitUrl parameter in URL");
            return;
        }

        if (params.submitMethod != null) {
            this.submitMethod = params.submitMethod;
        }

        if (params.successStatus != null) {
            this.successStatus = Number(params.successStatus);
        }

        this.submitUrl = params.submitUrl;
        this.schemaName = params.schemaName;

        if (params.value !== null) {
            this.value = params.value;
        }

        if (typeof params.valueUrl !== undefined) {
            this.valueUrl = params.valueUrl;
        }

        if (this.form == null) {
            this.setFormElement($('#schema-form'))
        }

        console.log("FORM", this)

        this.render()
    }

    JSONForm.prototype.render = function () {
        if (this.form === null) {
            this.error("Missing destination form element, did you call setFormElement?")
            return
        }

        var self = this

        if (this.schema === undefined) {
            var schemaUrl = this.schemaName + "-schema.json"

            send(schemaUrl, "GET", null, 200, function (resp) {
                console.log("SCHEMA RESP", resp)

                self.schema = JSON.parse(resp.responseText);

                if (self.title !== null && $(self.title).text() === '') {
                    $(self.title).text(self.schema.schema.title)
                    document.title = self.schema.schema.title
                }

                self.render()
            }, function (x) {
                self.error("Failed to load schema using URL:<br /><code>" + schemaUrl + "</code><br />Response:<br /><code>" + x.responseText + "</code>")
            })

            return
        }

        if (this.value === undefined && this.valueUrl !== undefined && this.valueUrl !== '') {
            send(this.valueUrl, "GET", null, 200, function (resp) {
                self.value = JSON.parse(resp.responseText);

                self.render()
            }, function (x) {
                self.error("Failed to load value using URL:<br /><code>" + self.valueUrl + "</code><br />Response:<br /><code>" + x.responseText + "</code>")
            })

            return
        }


        console.log("Rendering form")

        var formConf = {
            schema: this.schema.schema,
            form: this.schema.form,
            onSubmit: function (errors, values) {
                self.result.html('')

                console.log("VALUES", values);
                console.log("ERRORS", errors);

                if (errors) {
                    console.log(errors)
                    return;
                }

                if (self.beforeSubmit) {
                    self.beforeSubmit(values)
                }

                if (self.submitUrl && self.submitMethod) {
                    send(self.submitUrl, self.submitMethod, values, self.successStatus, self.success, self.fail)
                }
            }
        }

        if (typeof this.value !== undefined) {
            formConf.value = this.value
        }

        $(this.form).jsonForm(formConf);
    }

    /**
     * @param {Element} title - Title HTML element.
     */
    JSONForm.prototype.setTitleElement = function (title) {
        this.title = title;
    }


    /**
     * @param {Element} form - Form HTML element.
     */
    JSONForm.prototype.setFormElement = function (form) {
        this.form = form;
    }

    /**
     * @param {Element} result - Result HTML element.
     */
    JSONForm.prototype.setResultElement = function (result) {
        this.result = result;
    }

    /**
     * Get URL query params as a map.
     * @return {Object}
     */
    JSONForm.prototype.queryParams = function () {
        var query = window.location.search.substring(1);

        return query.split('&').reduce(function (res, item) {
            var parts = item.split('=');
            res[parts[0]] = decodeURIComponent(parts[1]);
            return res;
        }, {})
    }

    /**
     * Set schema name.
     * @param {string} name
     */
    JSONForm.prototype.setSchemaName = function (name) {
        this.schemaName = name;
    }

    /**
     * Set value source with a URL.
     * @param {string} url
     */
    JSONForm.prototype.setValueUrl = function (url) {
        this.valueUrl = url;
    }

    JSONForm.prototype.setSubmitURL = function (method, url) {
        this.submitMethod = method;
        this.submitUrl = url;
    }


    /**
     * @callback RawCallback
     * @param {XMLHttpRequest} value
     */

    /**
     * @callback HTMLCallback
     * @param {String} value
     */

    /**
     * @callback JSONCallback
     * @param {Object} value
     */

    /**
     *
     * @param {String} url
     * @param {String} method
     * @param {Object} bodyValues
     * @param {Number} successStatus
     * @param {RawCallback} successCallback
     * @param {RawCallback} failCallback
     */
    function send(url, method, bodyValues, successStatus, successCallback, failCallback) {
        var x = new XMLHttpRequest();
        x.onreadystatechange = function () {
            if (x.readyState !== XMLHttpRequest.DONE) {
                return;
            }

            console.log("request finished with status", x.status, "expected status", successStatus)

            if (!successStatus) {
                if (typeof (successCallback) === 'function') {
                    successCallback(x);
                }

                return
            }


            switch (x.status) {
                case successStatus:
                    if (typeof (successCallback) === 'function') {
                        successCallback(x);
                    }
                    break;
                default:
                    if (typeof (failCallback) === 'function') {
                        failCallback(x)
                    } else {
                        throw {err: 'unexpected response', data: x};
                    }
            }
        };


        x.open(method, url, true);
        if (typeof bodyValues !== null) {
            x.setRequestHeader("Content-Type", "application/json; charset=utf-8");
            x.send(JSON.stringify(bodyValues));
            return;
        }

        x.send();
    }

    window.JSONForm = JSONForm;
})();