const axios = require("axios");

module.exports = async () => {
    const {image, version} = process.env;
    const resp = await axios.get(`https://quay.io/api/v1/repository/mongodb/${image}`);
    return resp.data.tags[version] !== undefined;
};
