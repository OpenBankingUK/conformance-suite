exports.replayMiddleware = (req, res) => {
  try {
    const response = res;
    return res
      .status(response.statusCode)
      .set(response.headers)
      .json(response.body);
  } catch (e) {
    return res.json(e);
  }
};
