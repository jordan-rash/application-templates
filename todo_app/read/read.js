(_, payload) => {
  try {
    const jsonString = String.fromCharCode.apply(null, payload);
    const todo = JSON.parse(jsonString);

    if (!todo.id) {
      return {
        keys: this.hostServices.kv.keys(),
        status: "success"
      }
    } else {
      return {
        todo: String.fromCharCode(...this.hostServices.kv.get(todo.id)),
        status: "success"
      }
    }
  } catch (error) {
    return {
      status: "failed",
      error: error
    }
  }
};
