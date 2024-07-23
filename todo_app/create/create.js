(_, payload) => {
  try {
    const jsonString = String.fromCharCode.apply(null, payload);
    const todo = JSON.parse(jsonString);

    this.hostServices.kv.set(todo.id, payload);
    return {
      id: todo.id,
      status: "success"
    }
  } catch (error) {
    return {
      status: "failed",
      error: error
    }
  }
};
